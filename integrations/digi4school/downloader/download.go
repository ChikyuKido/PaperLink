package downloader

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"paperlink_d4s/downloader/helper"
	"paperlink_d4s/downloader/types"
	"paperlink_d4s/structs"
	"sort"
)

func DownloadBook(book *structs.Book, outputPath string, digi4sCookie string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	client.Jar, _ = cookiejar.New(nil)
	uri, _ := url.Parse("https://digi4school.at")
	client.Jar.SetCookies(uri, []*http.Cookie{
		{Name: "digi4s", Value: digi4sCookie},
	})

	data, _, lastURL, location, err := helper.GetLastLTI(client, book.DataCode)
	if err != nil {
		return nil
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

	var files []string
	if book.EbookPlus {
		if lastURL == "https://a.hpthek.at/lti" {
			files, err = types.DownloadD4sBook(client, tmp, location)
			if err != nil {
				return fmt.Errorf("failed to download book: %w", err)
			}
		} else if lastURL == "https://mein.westermann.de/auth/gateway/d4s" {
			files, err = types.DownloadBiboxBook(client, location, tmp)
			if err != nil {
				return fmt.Errorf("failed to download book: %w", err)
			}
		} else if lastURL == "https://service.helbling.com/ebookplus" {
			files, err = types.DownloadHeblingBook(client, data, tmp)
			if err != nil {
				return fmt.Errorf("failed to download book: %w", err)
			}
		} else {
			return fmt.Errorf("book source not supported")
		}
	} else {
		files, err = types.DownloadD4sBook(client, tmp, location)
		if err != nil {
			return fmt.Errorf("failed to download book: %w", err)
		}
	}

	sort.Strings(files)
	err = api.MergeCreateFile(files, outputPath, false, nil)
	if err != nil {
		return fmt.Errorf("failed to write pvf: %w", err)
	}
	return nil
}
