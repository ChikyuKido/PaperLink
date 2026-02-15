package ptf

import (
	"encoding/binary"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io"
	"os"
	"os/exec"
	"paperlink/util"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var log = util.GroupLog("PTF")
var thumbnailWorkerCount = 8

type thumbnailJob struct {
	worker int
	start  int
	end    int
	dir    string
}

func WriteThumbnailPTFFromPDF(inputFile string) (string, error) {
	pageCount, err := api.PageCountFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read pdf page count: %w", err)
	}
	if pageCount <= 0 {
		return "", fmt.Errorf("pdf has no pages")
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "pvf_thumb_*")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}

	workerCount := thumbnailWorkerCount
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > pageCount {
		workerCount = pageCount
	}
	jobs := splitThumbnailJobs(tempDir, pageCount, workerCount)

	if err := runThumbnailJobs(inputFile, jobs); err != nil {
		return "", err
	}
	pagePaths, err := collectJobPagePaths(jobs)
	if err != nil {
		return "", err
	}
	if len(pagePaths) != pageCount {
		return "", fmt.Errorf("thumbnail page count mismatch: expected=%d got=%d", pageCount, len(pagePaths))
	}

	outputFilePath := fmt.Sprintf("%s/output_thumb.ptf", tempDir)
	if err := writePTFByPagePaths(pagePaths, outputFilePath); err != nil {
		return "", err
	}
	return outputFilePath, nil
}

func splitThumbnailJobs(tempDir string, pageCount, workerCount int) []thumbnailJob {
	chunkSize := (pageCount + workerCount - 1) / workerCount
	jobs := make([]thumbnailJob, 0, workerCount)
	for i := 0; i < workerCount; i++ {
		start := i*chunkSize + 1
		if start > pageCount {
			break
		}
		end := start + chunkSize - 1
		if end > pageCount {
			end = pageCount
		}
		jobs = append(jobs, thumbnailJob{
			worker: i + 1,
			start:  start,
			end:    end,
			dir:    filepath.Join(tempDir, fmt.Sprintf("t%d", i+1)),
		})
	}
	return jobs
}

func runThumbnailJobs(inputFile string, jobs []thumbnailJob) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(jobs))

	for _, job := range jobs {
		job := job
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runThumbnailJob(inputFile, job); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func runThumbnailJob(inputFile string, job thumbnailJob) error {
	if err := os.MkdirAll(job.dir, 0o755); err != nil {
		return fmt.Errorf("failed to create thumbnail worker dir %s: %w", job.dir, err)
	}

	thumbnailPattern := filepath.Join(job.dir, "thumb_%05d.png")
	cmd := exec.Command(
		"gs",
		"-sDEVICE=pnggray",
		"-r100",
		"-dDownScaleFactor=4",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		fmt.Sprintf("-dFirstPage=%d", job.start),
		fmt.Sprintf("-dLastPage=%d", job.end),
		"-dBATCH",
		"-dNOPAUSE",
		"-sOutputFile="+thumbnailPattern,
		inputFile,
	)

	if out, cmdErr := cmd.CombinedOutput(); cmdErr != nil {
		return fmt.Errorf("ghostscript worker t%d failed for pages %d-%d: %w: %s", job.worker, job.start, job.end, cmdErr, strings.TrimSpace(string(out)))
	}
	return nil
}

func collectJobPagePaths(jobs []thumbnailJob) ([]string, error) {
	pagePaths := make([]string, 0)
	for _, job := range jobs {
		files, err := os.ReadDir(job.dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read thumbnail worker dir %s: %w", job.dir, err)
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".png" {
				continue
			}
			pagePaths = append(pagePaths, filepath.Join(job.dir, file.Name()))
		}
	}
	return pagePaths, nil
}

func writePTFByPagePaths(pagePaths []string, outputFilePath string) error {

	pageSizes := make([]uint64, 0, len(pagePaths))
	for _, pagePath := range pagePaths {
		info, err := os.Stat(pagePath)
		if err != nil {
			return fmt.Errorf("failed to stat split page %s: %w", filepath.Base(pagePath), err)
		}
		if info.Size() < 0 {
			return fmt.Errorf("invalid split page size for %s", filepath.Base(pagePath))
		}
		pageSizes = append(pageSizes, uint64(info.Size()))
	}

	pageCount := uint64(len(pagePaths))
	mapSize := uint64(8) + pageCount*indexEntrySize
	entries := make([]pageEntry, pageCount)
	nextOffset := uint64(headerFixedSize) + mapSize
	for i := uint64(0); i < pageCount; i++ {
		entries[i] = pageEntry{
			Offset: nextOffset,
			Size:   pageSizes[i],
		}
		nextOffset += pageSizes[i]
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(fileMagic[:]); err != nil {
		return fmt.Errorf("failed to write ptf magic: %w", err)
	}
	if _, err := outputFile.Write([]byte{fileVersionIndexed}); err != nil {
		return fmt.Errorf("failed to write ptf version: %w", err)
	}
	if err := binary.Write(outputFile, binary.LittleEndian, mapSize); err != nil {
		return fmt.Errorf("failed to write ptf map size: %w", err)
	}
	if err := binary.Write(outputFile, binary.LittleEndian, pageCount); err != nil {
		return fmt.Errorf("failed to write ptf page count: %w", err)
	}
	for i := uint64(0); i < pageCount; i++ {
		if err := binary.Write(outputFile, binary.LittleEndian, entries[i].Offset); err != nil {
			return fmt.Errorf("failed to write ptf page offset for page %d: %w", i, err)
		}
		if err := binary.Write(outputFile, binary.LittleEndian, entries[i].Size); err != nil {
			return fmt.Errorf("failed to write ptf page size for page %d: %w", i, err)
		}
	}

	pageCountInt := 0
	var totalPayload uint64
	for i, pagePath := range pagePaths {
		in, err := os.Open(pagePath)
		if err != nil {
			return fmt.Errorf("failed to open split page %s: %w", filepath.Base(pagePath), err)
		}

		written, err := io.Copy(outputFile, in)
		closeErr := in.Close()
		if err != nil {
			return fmt.Errorf("failed to write page data for %s: %w", filepath.Base(pagePath), err)
		}
		if closeErr != nil {
			return fmt.Errorf("failed to close split page %s: %w", filepath.Base(pagePath), closeErr)
		}
		if uint64(written) != pageSizes[i] {
			return fmt.Errorf("failed to write full page data for %s", filepath.Base(pagePath))
		}
		if err := os.Remove(pagePath); err != nil {
			return fmt.Errorf("failed to remove page %s: %w", filepath.Base(pagePath), err)
		}
		pageCountInt++
		totalPayload += pageSizes[i]
	}

	return nil
}
