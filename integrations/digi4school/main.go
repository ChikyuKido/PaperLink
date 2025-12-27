package main

import (
	"encoding/json"
	"fmt"
	"os"
	"paperlink_d4s/client"
	"strings"
)

type Status struct {
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
	BookName  string `json:"book_name"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: downloader list|download ...")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "list":
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "usage: downloader list <username> <password>")
			os.Exit(1)
		}
		username, password := os.Args[2], os.Args[3]
		if err := listBooks(username, password); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

	case "download":
		if len(os.Args) != 5 {
			fmt.Fprintln(os.Stderr, "usage: downloader download <id=path,...> <username> <password>")
			os.Exit(1)
		}
		mappingArg, username, password := os.Args[2], os.Args[3], os.Args[4]

		idPathMap, err := parseIDPathMapping(mappingArg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error parsing id/path mapping:", err)
			os.Exit(1)
		}

		if err := downloadBooks(idPathMap, username, password); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "unknown command:", cmd)
		fmt.Fprintln(os.Stderr, "available commands: list, download")
		os.Exit(1)
	}
}

func listBooks(username, password string) error {
	c := client.NewDigi4SClient(username, password)
	defer c.Logout()

	if err := c.Login(); err != nil {
		return err
	}

	books, err := c.GetBooks()
	if err != nil {
		return err
	}

	b, err := json.Marshal(books)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func downloadBooks(idPathMap map[string]string, username, password string) error {
	c := client.NewDigi4SClient(username, password)
	defer c.Logout()

	if err := c.Login(); err != nil {
		return err
	}

	books, err := c.GetBooks()
	if err != nil {
		return err
	}

	totalBooks := len(idPathMap)
	completed := 0

	for id, path := range idPathMap {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		found := false
		for _, b := range books {
			if b.DataId == id {
				found = true
				status := Status{
					Total:     totalBooks,
					Completed: completed,
					BookName:  b.Name,
				}
				statusJSON, _ := json.Marshal(status)
				fmt.Println(string(statusJSON))

				if err := c.DownloadBook(&b, path); err != nil {
					return fmt.Errorf("failed to download %s: %w", id, err)
				}
				completed++
				break
			}
		}

		if !found {
			return fmt.Errorf("book with id %s not found", id)
		}
	}
	status := Status{
		Total:     totalBooks,
		Completed: completed,
		BookName:  "",
	}
	statusJSON, _ := json.Marshal(status)
	fmt.Println(string(statusJSON))
	return nil
}

func parseIDPathMapping(s string) (map[string]string, error) {
	result := make(map[string]string)
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid segment %q (expected id=path)", part)
		}

		id := strings.TrimSpace(kv[0])
		path := strings.TrimSpace(kv[1])
		if id == "" || path == "" {
			return nil, fmt.Errorf("invalid id or path in segment %q", part)
		}

		result[id] = path
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid id=path pairs found")
	}

	return result, nil
}
