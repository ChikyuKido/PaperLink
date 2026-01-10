package helper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetLastLTI(c *http.Client, bookCode string) (string, map[string]string, string, string, error) {
	data, err := getOauthMap(c, bookCode)
	if err != nil {
		return "", nil, "", "", err
	}

	params := extractParams(data)
	lastParams := params
	actionURL := getActionURL(data)
	location := ""
	for {
		lastParams = params
		data, location, err = lti(c, actionURL, params)
		if err != nil {
			return "", nil, "", "", err
		}
		params = extractParams(data)
		if len(params) == 0 {
			break
		}
		actionURL = getActionURL(data)
	}
	return data, lastParams, actionURL, location, nil
}

func lti(c *http.Client, actionURL string, params map[string]string) (string, string, error) {
	baseUrl := actionURL
	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, url.QueryEscape(key)+"="+url.QueryEscape(value))
	}

	encodedFormData := strings.Join(queryParams, "&")
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBufferString(encodedFormData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://digi4school.at")
	req.Header.Set("Priority", "u=0, i")
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("could not load: %v", err)
	}
	location := resp.Header.Get("location")

	return string(body), location, nil
}

func getOauthMap(c *http.Client, bookCode string) (string, error) {
	baseUrl := "https://digi4school.at/ebook/" + bookCode

	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Priority", "u=0, i")
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}

func extractParams(formHTML string) map[string]string {
	params := make(map[string]string)
	re := regexp.MustCompile(`<input\s+name='([^']+)' value='([^']*)'>`)
	matches := re.FindAllStringSubmatch(formHTML, -1)
	for _, match := range matches {
		params[match[1]] = match[2]
	}
	return params
}

func getActionURL(formHTML string) string {
	re := regexp.MustCompile(`action=['"]([^'"]+)['"]`)
	matches := re.FindStringSubmatch(formHTML)
	return matches[1]
}
