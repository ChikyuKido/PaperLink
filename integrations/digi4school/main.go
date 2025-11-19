package main

import (
	"encoding/json"
	"fmt"
	"os"
	"paperlink_d4s/client"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  downloader list <username> <password>")
		fmt.Println("  downloader download <id1,id2,...> <username> <password>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "list":
		if len(os.Args) != 4 {
			fmt.Println("usage: downloader list <username> <password>")
			os.Exit(1)
		}
		username := os.Args[2]
		password := os.Args[3]

		if err := handleList(username, password); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

	case "download":
		if len(os.Args) != 5 {
			fmt.Println("usage: downloader download <id1=path1,id2=path2,...> <username> <password>")
			os.Exit(1)
		}
		mappingArg := os.Args[2]
		username := os.Args[3]
		password := os.Args[4]

		idPathMap, err := parseIDPathMapping(mappingArg)
		if err != nil {
			fmt.Println("error parsing id/path mapping:", err)
			os.Exit(1)
		}

		if err := handleDownload(idPathMap, username, password); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown command:", cmd)
		fmt.Println("available commands: list, download")
		os.Exit(1)
	}
}

func handleList(username, password string) error {
	c := client.NewDigi4SClient(username, password)
	err := c.Login()
	if err != nil {
		return err
	}
	books, _ := c.GetBooks()
	b, err := json.Marshal(books)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	c.Logout()
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

func handleDownload(idPathMap map[string]string, username, password string) error {
	c := client.NewDigi4SClient(username, password)
	err := c.Login()
	if err != nil {
		return err
	}
	for id, path := range idPathMap {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		books, _ := c.GetBooks()
		for _, b := range books {
			if b.DataId == id {
				err := c.DownloadBook(&b, path)
				if err != nil {
					return err
				}
			}
		}
	}

	c.Logout()
	return nil
}
