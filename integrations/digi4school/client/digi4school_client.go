package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"paperlink_d4s/downloader"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Digi4SchoolClient struct {
	Username string
	Password string
	Client   *http.Client
}

type BookCookies struct {
	Digi4Bname  string
	Digi4Bvalue string
	Digi4Pname  string
	Digi4Pvalue string
	SubPath     string
}

type Book struct {
	Name     string `json:"name"`
	DataCode string `json:"dataCode"`
	DataId   string `json:"dataId"`
}

func NewDigi4SClient(username, password string) *Digi4SchoolClient {
	transport := http.DefaultTransport
	jar, _ := cookiejar.New(nil)
	return &Digi4SchoolClient{
		Username: username,
		Password: password,
		Client: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		},
	}
}
func (c *Digi4SchoolClient) Login() error {
	baseUrl := "https://digi4school.at/br/xhr/login"

	payload := url.Values{}
	payload.Set("email", c.Username)
	payload.Set("password", c.Password)

	headers := map[string]string{
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0",
		"Accept":           "text/plain, */*; q=0.01",
		"Accept-Language":  "en-US,en;q=0.5",
		"Referer":          "https://digi4school.at/",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://digi4school.at",
	}

	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, cookie := range c.Client.Jar.Cookies(req.URL) {
		if cookie.Name == "digi4s" {
			return nil
		}
	}
	return fmt.Errorf("login failed")
}
func (c *Digi4SchoolClient) Logout() error {
	baseUrl := "https://digi4school.at/br/logout"

	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Referer":         "https://digi4school.at/",
	}

	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
func (c *Digi4SchoolClient) GetBooks() ([]Book, error) {
	baseURL := "https://digi4school.at/br/xhr/v2/synch"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	type apiResponse struct {
		Books []struct {
			ID        int    `json:"id"`
			Title     string `json:"title"`
			Code      string `json:"code"`
			Publisher string `json:"publisher"`
		} `json:"books"`
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	books := make([]Book, 0, len(result.Books))
	for _, b := range result.Books {
		books = append(books, Book{
			Name:     b.Title,
			DataCode: b.Code,
			DataId:   fmt.Sprintf("%d", b.ID),
		})
	}

	return books, nil
}
func (c *Digi4SchoolClient) DownloadBook(book *Book, outputPath string) error {
	bookCookies, err := c.getBookCookies(book.DataCode)
	if err != nil {
		return fmt.Errorf("could not get bookCookies: %w", err)
	}

	tmp, err := os.MkdirTemp("", "bookdl_*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	current, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current dir: %w", err)
	}
	defer os.Chdir(current)

	if err := os.Chdir(tmp); err != nil {
		return fmt.Errorf("failed to change dir: %w", err)
	}

	digi4bCookie := &http.Cookie{Name: bookCookies.Digi4Bname, Value: bookCookies.Digi4Bvalue}
	digi4pCookie := &http.Cookie{Name: bookCookies.Digi4Pname, Value: bookCookies.Digi4Pvalue}
	digi4sCookie := &http.Cookie{Name: "digi4s", Value: c.getCurrentDigi4sCookie()}
	cookies := []*http.Cookie{digi4sCookie, digi4pCookie, digi4bCookie}

	page := 1
	jobs := make(chan string, 1000)
	results := make(chan string, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 12; i++ {
		wg.Add(1)
		go svgWorker(jobs, results, tmp, &wg)
	}

	for {
		var baseUrl string
		if bookCookies.SubPath != "" {
			baseUrl = fmt.Sprintf("https://a.digi4school.at/ebook/%s/%s", book.DataId, bookCookies.SubPath)
		} else {
			baseUrl = fmt.Sprintf("https://a.digi4school.at/ebook/%s", book.DataId)
		}

		name, err := downloader.DownloadOnePage(fmt.Sprintf("%s/%d.svg", baseUrl, page), cookies)
		if name != "" {
			jobs <- name
		}

		if err != nil {
			break
		}

		page++
	}

	close(jobs)
	wg.Wait()
	close(results)

	if page-1 != len(results) {
		return fmt.Errorf("unexpected page count: got %d, expected %d", len(results), page-1)
	}

	var outputPDFs []string
	for outputPDF := range results {
		outputPDFs = append(outputPDFs, outputPDF)
	}
	sort.Strings(outputPDFs)

	err = api.MergeCreateFile(outputPDFs, outputPath, false, &model.Configuration{
		Optimize: true,
	})
	if err != nil {
		return fmt.Errorf("failed to merge PDFs: %w", err)
	}

	return nil
}

func svgWorker(jobs <-chan string, results chan<- string, tempDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	for svgFile := range jobs {
		outputPDF := filepath.Join(tempDir, strings.TrimSuffix(svgFile, ".svg")+".pdf")
		cmd := exec.Command(
			"rsvg-convert",
			"-f", "pdf",
			"-o", outputPDF,
			filepath.Join(tempDir, svgFile),
		)
		if err := cmd.Run(); err != nil {
			continue
		}
		results <- outputPDF
	}
}

func (c *Digi4SchoolClient) getBookCookies(bookId string) (BookCookies, error) {
	oauthMap, err := c.getOauthMap(bookId)
	if err != nil {
		return BookCookies{}, fmt.Errorf("could not refresh digi4s cookie: %v", err)
	}
	oauthMap2, _ := c.lti1Request(oauthMap)
	finishedCookies, _ := c.lti2Request(oauthMap2)
	return finishedCookies, nil
}
func (c *Digi4SchoolClient) lti1Request(params map[string]string) (map[string]string, error) {
	baseUrl := "https://kat.digi4school.at/lti"
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

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]string{}, fmt.Errorf("could not load: %v", err)
	}
	return extractParams(string(body)), nil
}
func (c *Digi4SchoolClient) lti2Request(params map[string]string) (BookCookies, error) {
	baseUrl := "https://a.digi4school.at/lti"
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
	req.Header.Set("Origin", "https://kat.digi4school.at")
	req.Header.Set("Priority", "u=0, i")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	finishedCookies := BookCookies{}

	finishedCookies.SubPath = c.checkSubPath(resp.Header.Get("Location"))

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "digi4b" {
			finishedCookies.Digi4Bvalue = cookie.Value
			finishedCookies.Digi4Bname = cookie.Name
		}
		if cookie.Name == "digi4p" {
			finishedCookies.Digi4Pvalue = cookie.Value
			finishedCookies.Digi4Pname = cookie.Name
		}
	}
	if finishedCookies.Digi4Bvalue == "" || finishedCookies.Digi4Pvalue == "" {
		//error handling
	}
	return finishedCookies, nil
}
func (c *Digi4SchoolClient) checkSubPath(url string) string {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://kat.digi4school.at")
	req.Header.Set("Priority", "u=0, i")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "sbnr") {
		return "1"
	}

	return ""
}
func (c *Digi4SchoolClient) getOauthMap(buchId string) (map[string]string, error) {
	baseUrl := "https://digi4school.at/ebook/" + buchId

	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Priority", "u=0, i")

	resp, err := c.Client.Do(req)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return extractParams(string(body)), nil
}
func (c *Digi4SchoolClient) getCurrentDigi4sCookie() string {
	uri, _ := url.Parse("https://a.digi4school.at")
	for _, cookie := range c.Client.Jar.Cookies(uri) {
		if cookie.Name == "digi4s" {
			return cookie.Value
		}
	}
	return ""
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
