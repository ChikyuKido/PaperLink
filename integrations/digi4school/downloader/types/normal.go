package types

import (
	"fmt"
	"net/http"
	"os"
	"paperlink_d4s/downloader/helper"
)

func DownloadD4sBook(c *http.Client, baseURL string, downloadPath string) ([]string, error) {
	subPath, continuingSubPath := withSubPath(c, baseURL)
	page := 1
	err := os.Chdir(downloadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to chdir to %s: %w", downloadPath, err)
	}
	files := make([]string, 0)
	for {
		downloadURL := fmt.Sprintf("%s/%d.svg", baseURL, page)
		if continuingSubPath {
			downloadURL = fmt.Sprintf("%s/%d/%d.svg", baseURL, page, page)
		} else if subPath {
			downloadURL = fmt.Sprintf("%s/1/%d.svg", baseURL, page)
		}
		filename, endReached, err := helper.DownloadOnePage(downloadURL, c, subPath)
		if err != nil {
			return nil, fmt.Errorf("failed to download page: %w", err)
		}
		if endReached {
			break
		}
		outputPDF, err := helper.ConvertAndCompressSVG(downloadPath, filename)
		if err != nil {
			continue
		}
		files = append(files, outputPDF)
		page++
	}

	return files, nil
}

func withSubPath(c *http.Client, baseURL string) (bool, bool) {
	req, err := http.NewRequest("GET", baseURL+"/1/1.svg", nil)
	resp, err := c.Do(req)
	if err != nil {
		return false, false
	}
	subPath := false
	if resp.StatusCode != 404 {
		subPath = true
	}
	req, err = http.NewRequest("GET", baseURL+"/2/2.svg", nil)
	resp, err = c.Do(req)
	continuingSubPath := false
	if err != nil {
		return false, false
	}
	if resp.StatusCode != 404 {
		continuingSubPath = true
	}
	return subPath, continuingSubPath
}
