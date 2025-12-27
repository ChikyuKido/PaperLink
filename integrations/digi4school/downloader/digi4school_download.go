package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func downloadEmbeddedAsset(url string, matches [][]string, cookies []*http.Cookie) error {
	trimmedURL := url[:strings.LastIndex(url, "/")+1]
	for _, match := range matches {
		if len(match) > 1 {
			if err := downloadFile(fmt.Sprintf(trimmedURL+match[1]), cookies); err != nil {
				return fmt.Errorf("failed to download embedded asset %s: %w", match[1], err)
			}
		}
	}
	return nil
}

func downloadFile(url string, cookies []*http.Cookie) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", url, err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get file %s: %w", url, err)
	}
	defer resp.Body.Close()

	dirname := getDirName(url)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirname, err)
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

func DownloadOnePage(url string, cookies []*http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", url, err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get file %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("ERR 404 - %s", url)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body for %s: %w", url, err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "image") {
		matches := checkForEmbeddedImages(bodyString)
		if len(matches) > 0 {
			if err := downloadEmbeddedAsset(url, matches, cookies); err != nil {
				return "", fmt.Errorf("failed to download embedded images for %s: %w", url, err)
			}
		}
	}

	filename := path.Base(url)
	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid file name %s", filename)
	}
	number := strings.Repeat("0", 5-len(parts[0])) + parts[0]
	filename = number + "." + parts[1]

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if _, err := io.WriteString(file, bodyString); err != nil {
		return "", fmt.Errorf("failed to write content to %s: %w", filename, err)
	}

	return filename, nil
}
