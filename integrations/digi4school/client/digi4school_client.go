package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"paperlink_d4s/structs"
	"strings"
	"time"
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
func (c *Digi4SchoolClient) GetBooks() ([]structs.Book, error) {
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
			EbookPlus int    `json:"ebook_plus"`
			Expiry    string `json:"expiry"`
		} `json:"books"`
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	books := make([]structs.Book, 0, len(result.Books))
	for _, b := range result.Books {
		if ok := isExpired(b.Expiry); ok {
			continue
		}
		ebookPlus := false
		if b.EbookPlus == 1 {
			ebookPlus = true
		}
		books = append(books, structs.Book{
			Name:      b.Title,
			DataCode:  b.Code,
			DataId:    fmt.Sprintf("%d", b.ID),
			EbookPlus: ebookPlus,
			Publisher: b.Publisher,
		})
	}

	return books, nil
}
func (c *Digi4SchoolClient) GetCurrentDigi4sCookie() string {
	uri, _ := url.Parse("https://a.digi4school.at")
	for _, cookie := range c.Client.Jar.Cookies(uri) {
		if cookie.Name == "digi4s" {
			return cookie.Value
		}
	}
	return ""
}

func isExpired(dateStr string) bool {
	const layout = "02.01.2006"

	t, err := time.ParseInLocation(layout, dateStr, time.Local)
	if err != nil {
		return true
	}

	expired := t.Before(time.Now())
	return expired
}
