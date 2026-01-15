package types

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"paperlink_d4s/downloader/helper"
	"regexp"
	"strconv"
	"strings"
)

type Image struct {
	ID       int    `json:"id"`
	Removed  bool   `json:"removed"`
	Version  int    `json:"version"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	PageID   int    `json:"pageId"`
	MD5Sum   string `json:"md5sum"`
	FileSize int    `json:"filesize"`
}

type Page struct {
	ID              int     `json:"id"`
	Removed         bool    `json:"removed"`
	Version         int     `json:"version"`
	Name            string  `json:"name"`
	InternalPageNum int     `json:"internalPagenum"`
	BookID          int     `json:"bookId"`
	Demo            bool    `json:"demo"`
	Type            string  `json:"type"`
	AemDorisID      *string `json:"aemDorisID"`
	Images          []Image `json:"images"`
}

type BookPagesResponse struct {
	Pages []Page `json:"pages"`
}

func DownloadBiboxBook(c *http.Client, location string, downloadPath string) ([]string, error) {
	loginHint, id, ok := extractLoginInitParams(location)
	if !ok {
		return nil, fmt.Errorf("failed to extract loginHint")
	}
	token, err := getBiboxJWT(c, loginHint)
	if err != nil {
		return nil, err
	}
	pages, err := getBookPages(c, id, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get bibox pages: %w", err)
	}
	err = os.Chdir(downloadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to chdir: %w", err)
	}
	files := make([]string, 0)
	page := 1
	for _, p := range pages {
		images := p.Images
		if len(images) == 0 {
			return nil, fmt.Errorf("no images found for page %d", page)
		}
		downloadImage := images[0]
		// download the higher res version
		if len(images) == 2 {
			if images[1].FileSize > images[0].FileSize {
				downloadImage = images[1]
			}
		}

		number := strings.Repeat("0", 5-len(strconv.Itoa(page))) + strconv.Itoa(page)
		outputFile := number + ".png"
		resp, err := http.Get(downloadImage.URL)
		if err != nil {
			panic(err)
		}
		outFile, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		outFile.Close()

		pdf, err := helper.ConvertPNGtoPDF(downloadPath, outputFile, downloadImage.Width, downloadImage.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to convert png to pdf: %w", err)
		}
		files = append(files, pdf)
		page++
		fmt.Printf("PAGE_COUNT: %d\n", page)
	}
	return files, nil
}

func getBookPages(c *http.Client, bookID int, jwt string) ([]Page, error) {
	bookUrl := fmt.Sprintf("https://backend.bibox2.westermann.de/v1/api/sync/%d?materialtypes[]=default&materialtypes[]=addon", bookID)

	req, err := http.NewRequest("GET", bookUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Origin", "https://bibox2.westermann.de")
	req.Header.Set("Referer", "https://bibox2.westermann.de/")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", body)
	}

	var result BookPagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Pages, nil
}

func getBiboxJWT(c *http.Client, loginHint string) (string, error) {
	codeVerifier := randomString(64)
	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])

	state := randomString(16)

	q := url.Values{}
	q.Set("client_id", "Nvw0ZA8Z")
	q.Set("response_type", "code")
	q.Set("scope", "openid")
	q.Set("redirect_uri", "https://bibox2.westermann.de/login")
	q.Set("state", state)
	q.Set("code_challenge_method", "S256")
	q.Set("code_challenge", codeChallenge)
	q.Set("login_hint", loginHint)

	req, err := http.NewRequest(
		"GET",
		"https://mein.westermann.de/auth/login?"+q.Encode(),
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Referer", "https://bibox2.westermann.de/")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	code, ok := extractAuthCode(string(body))
	if !ok {
		return "", fmt.Errorf("failed to extract auth code")
	}

	payload := map[string]string{
		"code":          code,
		"code_verifier": codeVerifier,
		"redirect_uri":  "https://bibox2.westermann.de/login",
	}

	b, _ := json.Marshal(payload)

	tokenReq, err := http.NewRequest(
		"POST",
		"https://backend.bibox2.westermann.de/token",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return "", err
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenReq.Header.Set("Accept", "application/json")
	tokenReq.Header.Set("User-Agent", "Mozilla/5.0")

	tokenResp, err := c.Do(tokenReq)
	if err != nil {
		return "", err
	}
	defer tokenResp.Body.Close()

	tokenBody, _ := io.ReadAll(tokenResp.Body)
	var res struct {
		IDToken string `json:"id_token"`
	}

	if err := json.Unmarshal(tokenBody, &res); err != nil {
		return "", err
	}
	if res.IDToken == "" {
		return "", fmt.Errorf("no id_token in response")
	}

	return res.IDToken, nil

}
func extractAuthCode(input string) (string, bool) {
	re := regexp.MustCompile(`code=([^&"]+)`)
	m := re.FindStringSubmatch(input)
	if len(m) < 2 {
		return "", false
	}
	return m[1], true
}
func extractLoginInitParams(input string) (loginHint string, bookID int, ok bool) {
	re := regexp.MustCompile(`login_hint=([^&]+).*?target_link_uri=%2Fv2%2Fbook%2F(\d+)`)
	m := re.FindStringSubmatch(input)
	if len(m) != 3 {
		return "", 0, false
	}
	id, _ := strconv.Atoi(m[2])

	return m[1], id, true
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}
