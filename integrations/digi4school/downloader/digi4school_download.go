package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func downloadEmbeddedAsset(url string, matches [][]string, cookies []*http.Cookie) {
	trimmedURL := url[:strings.LastIndex(url, "/")+1]
	for _, match := range matches {
		if len(match) > 1 {
			downloadFile(fmt.Sprintf(trimmedURL+match[1]), cookies)
		}
	}
}

func downloadFile(url string, cookies []*http.Cookie) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicln("failed to create request")
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("failed to get file:", err)
	}
	defer resp.Body.Close()

	dirname := getDirName(url)

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}

	file, err := os.Create(dirname + path.Base(url))
	if err != nil {
		log.Panicln("failed to create file")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Panicln("failed to copy file content:", err)
	}

}

func DownloadOnePage(url string, cookies []*http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicln("failed to create request")
		return "", err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("failed to get file:", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("ERR 404 - %s", url)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
		log.Panicln("Error:", err)
	}

	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "image") {
		matches := checkForEmbeddedImages(bodyString)
		if len(matches) > 0 {
			downloadEmbeddedAsset(url, matches, cookies)
		}
	}
	filename := path.Base(url)
	var parts = strings.Split(filename, ".")
	var length = len(parts[0])
	var number = strings.Repeat("0", 5-length) + parts[0]
	filename = number + "." + parts[1]

	file, err := os.Create(filename)
	if err != nil {
		log.Panicln("failed to create file")
		return "", err
	}
	defer file.Close()

	_, err = io.WriteString(file, bodyString)
	if err != nil {
		log.Panicln("failed to copy file content:", err)
		return "", err
	}
	return filename, nil
}
