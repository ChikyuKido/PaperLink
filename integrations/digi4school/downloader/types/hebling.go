package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"paperlink_d4s/downloader/helper"
	"regexp"
)

func DownloadHeblingBook(c *http.Client, data string, downloadPath string) ([]string, error) {
	jwt, book := extractData(data)
	baseURL, err := getBookBaseURL(c, jwt, book)
	if err != nil {
		return nil, fmt.Errorf("error getting book url: %w", err)
	}
	page := 1
	err = os.Chdir(downloadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to chdir to %s: %w", downloadPath, err)
	}
	files := make([]string, 0)
	for {
		downloadURL := fmt.Sprintf("%s/pages/svg/%d.svg", baseURL, page)
		filename, endReached, err := helper.DownloadOnePage(downloadURL, c, false)
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

func getBookBaseURL(c *http.Client, jwt string, book string) (string, error) {
	url := fmt.Sprintf("https://service.helbling.com/api/productItems/%s/reference", book)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	var data struct {
		Content string `json:"content"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}
	return data.Content, nil
}
func extractData(html string) (jwt, lastPart string) {
	tokenRe := regexp.MustCompile(`window\.localStorage\.setItem\('d4s-token',\s*'([^']+)'`)
	hrefRe := regexp.MustCompile(`window\.location\.href\s*=\s*'([^']+)'`)
	lastPartRe := regexp.MustCompile(`#/ebook/[^/]+/([^']+)$`)
	jwtRe := regexp.MustCompile(`"access_token":"([^"]+)"`)

	var tokenJSON string
	if match := tokenRe.FindStringSubmatch(html); len(match) > 1 {
		tokenJSON = match[1]
	}
	if tokenJSON != "" {
		if match := jwtRe.FindStringSubmatch(tokenJSON); len(match) > 1 {
			jwt = match[1]
		}
	}
	if match := hrefRe.FindStringSubmatch(html); len(match) > 1 {
		href := match[1]
		if partMatch := lastPartRe.FindStringSubmatch(href); len(partMatch) > 1 {
			lastPart = partMatch[1]
		}
	}

	return
}
