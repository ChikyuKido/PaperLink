package d4s

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
	"os/exec"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/ptf"
	"paperlink/service/task"
	"paperlink/util"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StartSyncTask(accs []entity.Digi4SchoolAccount) (string, error) {
	control := &syncControl{}
	l, err := task.CreateNewTask("Digi4School Sync", control.Stop)
	if err != nil {
		return "", err
	}
	go syncAccounts(l, accs, control)
	return l.Task.ID, nil
}

func syncAccounts(l *task.TaskRunner, accs []entity.Digi4SchoolAccount, control *syncControl) {
	l.Info(fmt.Sprintf("Sync %d accounts", len(accs)))
	accountBooks := make([]Book, 0)
	for _, acc := range accs {
		if !l.IsRunning() || control.IsStopRequested() {
			l.Warn("task stopped by user")
			return
		}
		l.Info(fmt.Sprintf("Search books for account: %s", acc.Username))
		books, err := ListBooksForAccount(&acc)
		if err != nil {
			l.Err(fmt.Sprintf("Failed to list books for account: %s", acc.Username))
		}
		l.Info(fmt.Sprintf("Found %d books for account: %s", len(books), acc.Username))
		accountBooks = append(accountBooks, books...)
	}
	dbBooks, err := repo.Digi4SchoolBook.GetList()
	if err != nil {
		l.Critical(fmt.Sprintf("Failed to list books for account: %s", err.Error()))
		err := l.Fail()
		if err != nil {
			log.Error("Could not fail the running task")
		}
		return
	}
	l.Info(fmt.Sprintf("Found %d Books in %d accounts. %d Books are already in the db", len(accountBooks), len(accs), len(dbBooks)))
	unqiueAccountBooksMap := make(map[string]Book)
	for _, book := range accountBooks {
		unqiueAccountBooksMap[book.DataId] = book
	}
	uniqueAccountBooks := make([]Book, 0)
	for _, book := range unqiueAccountBooksMap {
		uniqueAccountBooks = append(uniqueAccountBooks, book)
	}
	dbBooksMap := make(map[string]entity.Digi4SchoolBook)
	for _, book := range dbBooks {
		dbBooksMap[book.BookID] = book
	}
	neededBooks := make([]Book, 0)
	for _, book := range uniqueAccountBooks {
		if _, ok := dbBooksMap[book.DataId]; !ok {
			neededBooks = append(neededBooks, book)
		}
	}
	l.Info(fmt.Sprintf("Found %d needed books. Start downloading", len(neededBooks)))
	err = downloadBooks(l, neededBooks, control)
	if control.IsStopRequested() {
		l.Warn("sync task stopped by user")
		return
	}
	if err != nil {
		err := l.Fail()
		if err != nil {
			log.Error("Failed to fail the sync task")
		}
		return
	}
	err = l.Complete()
	if err != nil {
		log.Error("Failed to complete the sync task")
	}
}
func downloadBooks(l *task.TaskRunner, books []Book, control *syncControl) error {
	copyBooks := slices.Clone(books)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(wd, "data", "d4s")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseDir, 0750)
		if err != nil {
			return fmt.Errorf("cannot create directory %s: %v", baseDir, err)
		}
	}
	control.SetRescanContext(baseDir, copyBooks)
	for len(books) > 0 {
		if !l.IsRunning() || control.IsStopRequested() {
			return nil
		}
		username := books[0].Account.Username
		sameAccountBooks := make([]Book, 0)
		i := 0
		for i < len(books) {
			if books[i].Account.Username == username {
				sameAccountBooks = append(sameAccountBooks, books[i])
				books = append(books[:i], books[i+1:]...)
			} else {
				i++
			}
		}
		if len(sameAccountBooks) == 0 {
			break
		}
		l.Info(fmt.Sprintf("Download %d books for account %s", len(sameAccountBooks), sameAccountBooks[0].Account.Username))
		var downloadIdString strings.Builder
		for _, book := range sameAccountBooks {
			if downloadIdString.Len() > 0 {
				downloadIdString.WriteString(",")
			}
			downloadIdString.WriteString(book.DataId)
			downloadIdString.WriteString("=")
			downloadIdString.WriteString(filepath.Join(baseDir, book.UUID+".pdf"))
		}
		acc := sameAccountBooks[0].Account
		cmd := exec.Command("./integrations/d4s", "download", downloadIdString.String(), acc.Username, acc.Password)
		control.SetCurrentCmd(cmd)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			control.ClearCurrentCmd(cmd)
			return fmt.Errorf("failed to start downloader command: %w", err)
		}
		go func() {
			scanner := bufio.NewScanner(stdout)
			lastText := ""
			start := time.Now()
			var lastPageTime time.Time
			var lastPage int

			for scanner.Scan() {
				t := scanner.Text()
				if strings.Contains(t, "PAGE_COUNT") && strings.Contains(lastText, "PAGE_COUNT") {
					parts := strings.Split(t, ": ")
					page, _ := strconv.Atoi(parts[1])
					now := time.Now()

					if lastPageTime.IsZero() {
						lastPageTime = now
					}

					elapsed := now.Sub(start)
					pagesSinceLast := page - lastPage
					secondsSinceLast := now.Sub(lastPageTime).Seconds()
					avgPerPage := elapsed.Seconds() / float64(page)
					pagesPerSecond := float64(pagesSinceLast) / secondsSinceLast

					l.ReplaceLastInfo(fmt.Sprintf(
						"Downloaded %d Pages | Avg time/page: %.2fs | Pages/sec: %.2f | Total elapsed: %s",
						page, avgPerPage, pagesPerSecond, elapsed.Truncate(1*time.Second),
					))

					lastPage = page
					lastPageTime = now
				} else {
					l.Info(t)
				}
				lastText = t
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				l.Err(fmt.Sprintf("Error occoured while downloading book: %s", scanner.Text()))
			}
		}()

		go func() {
			for {
				time.Sleep(180 * time.Second)
				if !l.IsRunning() || control.IsStopRequested() {
					return
				}
				err := rescanForDBInsert(baseDir, copyBooks)
				if err != nil {
					l.Err(fmt.Sprintf("Failed to rescan for books: %s", err.Error()))
				}
			}
		}()

		waitErr := cmd.Wait()
		control.ClearCurrentCmd(cmd)
		if waitErr != nil && !control.IsStopRequested() {
			l.Err(fmt.Sprintf("Downloader process exited with error: %v", waitErr))
		}
		err := rescanForDBInsert(baseDir, copyBooks)
		if err != nil {
			l.Err(fmt.Sprintf("Failed to rescan for books: %s", err.Error()))
		}
		if !l.IsRunning() || control.IsStopRequested() {
			return nil
		}
	}
	return nil
}

type syncControl struct {
	mu            sync.Mutex
	currentCmd    *exec.Cmd
	stopRequested bool
	rescanBaseDir string
	rescanBooks   []Book
}

func (s *syncControl) SetCurrentCmd(cmd *exec.Cmd) {
	s.mu.Lock()
	s.currentCmd = cmd
	s.mu.Unlock()
}

func (s *syncControl) ClearCurrentCmd(cmd *exec.Cmd) {
	s.mu.Lock()
	if s.currentCmd == cmd {
		s.currentCmd = nil
	}
	s.mu.Unlock()
}

func (s *syncControl) SetRescanContext(baseDir string, books []Book) {
	s.mu.Lock()
	s.rescanBaseDir = baseDir
	s.rescanBooks = slices.Clone(books)
	s.mu.Unlock()
}

func (s *syncControl) IsStopRequested() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopRequested
}

func (s *syncControl) Stop(l *task.TaskRunner) error {
	s.mu.Lock()
	s.stopRequested = true
	cmd := s.currentCmd
	baseDir := s.rescanBaseDir
	books := slices.Clone(s.rescanBooks)
	s.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		l.Warn("stopping digi4school downloader process")
		if err := cmd.Process.Kill(); err != nil {
			l.Err(fmt.Sprintf("failed to kill downloader process: %v", err))
		}
	}

	if baseDir != "" && len(books) > 0 {
		l.Info("running final rescan before stopping")
		if err := rescanForDBInsert(baseDir, books); err != nil {
			l.Err(fmt.Sprintf("final rescan failed: %v", err))
		}
	}
	return nil
}

func rescanForDBInsert(dir string, books []Book) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %v", dir, err)
	}
	for _, file := range files {
		for _, book := range books {
			// already in db
			if repo.Digi4SchoolBook.GetByUUID(book.UUID) != nil {
				continue
			}
			if file.Name() == book.UUID+".pdf" {
				fullPath := filepath.Join(dir, file.Name())
				info, statErr := os.Stat(fullPath)
				if statErr != nil {
					return fmt.Errorf("failed to stat file %s: %v", fullPath, statErr)
				}

				pageCount, err := api.PageCountFile(fullPath)
				if err != nil {
					return fmt.Errorf("failed to read page count file %s: %v", fullPath, err)
				}
				thumbPTFFile, err := ptf.WriteThumbnailPTFFromPDF(fullPath)
				if err != nil {
					return fmt.Errorf("failed to generate thumbnail ptf file %s: %v", fullPath, err)
				}
				thumbDst := strings.TrimSuffix(fullPath, filepath.Ext(fullPath)) + "_thumb.ptf"
				if err := util.CopyFile(thumbPTFFile, thumbDst); err != nil {
					_ = os.RemoveAll(filepath.Dir(thumbPTFFile))
					return fmt.Errorf("failed to store thumbnail ptf file %s: %v", thumbDst, err)
				}
				_ = os.RemoveAll(filepath.Dir(thumbPTFFile))

				fd := entity.FileDocument{
					UUID:  book.UUID,
					Path:  fullPath,
					Size:  uint64(info.Size()),
					Pages: uint64(pageCount),
				}
				_ = repo.FileDocument.Save(&fd)

				err = repo.Digi4SchoolBook.Save(&entity.Digi4SchoolBook{
					UUID:      book.UUID,
					BookName:  book.Name,
					BookID:    book.DataId,
					AccountID: book.Account.ID,
					FileUUID:  book.UUID,
				})
				if err != nil {
					return fmt.Errorf("failed to save book %s: %v", book.UUID, err)
				}
			}
		}
	}

	return nil
}
func ListBooksForAccount(acc *entity.Digi4SchoolAccount) ([]Book, error) {
	cmd := exec.Command("./integrations/d4s", "list", acc.Username, acc.Password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to execute list for user %s: %v, output: %s", acc.Username, err, string(output))
		return nil, err
	}

	outStr := string(output)
	outStr = strings.TrimSpace(outStr)
	var books []Book
	if err := json.Unmarshal([]byte(outStr), &books); err != nil {
		log.Printf("failed to unmarshal list for user %s: %v, output: %s", acc.Username, err, string(output))
		return nil, err
	}
	for i, _ := range books {
		books[i].Account = acc
		books[i].UUID = uuid.NewString()
	}
	return books, nil
}
