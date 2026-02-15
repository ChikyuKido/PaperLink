package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func downloadEmbeddedAsset(url string, matches [][]string, client *http.Client) error {
	trimmedURL := url[:strings.LastIndex(url, "/")+1]
	for _, match := range matches {
		if len(match) > 1 {
			if err := downloadFile(trimmedURL+match[1], client); err != nil {
				return fmt.Errorf("failed to download embedded asset %s: %w", match[1], err)
			}
		}
	}
	return nil
}

func downloadFile(url string, client *http.Client) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get file %s: %w", url, err)
	}
	defer resp.Body.Close()

	dirname := getDirName(url)
	if dirname != "" {
		if _, err := os.Stat(dirname); os.IsNotExist(err) {
			if err := os.MkdirAll(dirname, 0700); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dirname, err)
			}
		}
	}

	filePath := path.Join(dirname, path.Base(url))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to copy file content for %s: %w", filePath, err)
	}

	return nil
}

func DownloadOnePage(url string, client *http.Client, subPath bool) (string, bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", false, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("failed to get file %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", true, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to read response body for %s: %w", url, err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "image") {
		matches := checkForEmbeddedImages(bodyString)
		if len(matches) > 0 {
			if err := downloadEmbeddedAsset(url, matches, client); err != nil {
				return "", false, fmt.Errorf("failed to download embedded images for %s: %w", url, err)
			}
		}
	}

	filename := path.Base(url)
	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return "", false, fmt.Errorf("invalid file name %s", filename)
	}
	number := strings.Repeat("0", 5-len(parts[0])) + parts[0]
	filename = number + "." + parts[1]
	if subPath {
		pageDir := normalizePageDir(parts[0])
		filename = pageDir + "/" + filename
		err := os.MkdirAll(pageDir, 0700)
		if err != nil {
			return "", false, fmt.Errorf("failed to create directory %s: %w", pageDir, err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", false, fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if _, err := io.WriteString(file, bodyString); err != nil {
		return "", false, fmt.Errorf("failed to write content to %s: %w", filename, err)
	}

	return filename, false, nil
}

func normalizePageDir(raw string) string {
	page, err := strconv.Atoi(raw)
	if err != nil || page < 0 {
		return raw
	}
	return fmt.Sprintf("%04d", page)
}
