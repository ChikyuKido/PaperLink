package d4s

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/service/task"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func StartSyncTask(accs []entity.Digi4SchoolAccount) (string, error) {
	l, err := task.CreateNewTask("Digi4School Sync")
	if err != nil {
		return "", err
	}
	go syncAccounts(l, accs)
	return l.Task.ID, nil
}

func syncAccounts(l *task.TaskRunner, accs []entity.Digi4SchoolAccount) {
	l.Info(fmt.Sprintf("Sync %d accounts", len(accs)))
	accountBooks := make([]Book, 0)
	for _, acc := range accs {
		l.Info(fmt.Sprintf("Search books for account: %s", acc.Username))
		books, err := ListBooksForAccount(&acc)
		if err != nil {
			l.Err(fmt.Sprintf("Failed to list books for account: %s", acc.Username))
		}
		l.Info(fmt.Sprintf("Found %d books for account: %s", len(books), acc.Username))
		accountBooks = append(accountBooks, books...)
	}
	dbBooks, err := repo.Digi4SchoolBook.GetList()
	if err != nil {
		l.Critical(fmt.Sprintf("Failed to list books for account: %s", err.Error()))
		err := l.Fail()
		if err != nil {
			log.Error("Could not fail the running task")
		}
		return
	}
	l.Info(fmt.Sprintf("Found %d Books in %d accounts. %d Books are already in the db", len(accountBooks), len(accs), len(dbBooks)))
	unqiueAccountBooksMap := make(map[string]Book)
	for _, book := range accountBooks {
		unqiueAccountBooksMap[book.DataId] = book
	}
	uniqueAccountBooks := make([]Book, 0)
	for _, book := range unqiueAccountBooksMap {
		uniqueAccountBooks = append(uniqueAccountBooks, book)
	}
	dbBooksMap := make(map[string]entity.Digi4SchoolBook)
	for _, book := range dbBooks {
		dbBooksMap[book.BookID] = book
	}
	neededBooks := make([]Book, 0)
	for _, book := range uniqueAccountBooks {
		if _, ok := dbBooksMap[book.DataId]; !ok {
			neededBooks = append(neededBooks, book)
		}
	}
	neededBooks = []Book{neededBooks[0]}
	l.Info(fmt.Sprintf("Found %d needed books. Start downloading", len(neededBooks)))
	err = downloadBooks(l, neededBooks)
	if err != nil {
		err := l.Fail()
		if err != nil {
			log.Error("Failed to fail the sync task")
		}
	}
	err = l.Complete()
	if err != nil {
		log.Error("Failed to complete the sync task")
	}
}
func downloadBooks(l *task.TaskRunner, books []Book) error {
	copyBooks := slices.Clone(books)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(wd, "data", "d4s")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseDir, 0750)
		if err != nil {
			return fmt.Errorf("cannot create directory %s: %v", baseDir, err)
		}
	}
	for len(books) > 0 {
		username := books[0].Account.Username
		sameAccountBooks := make([]Book, 0)
		i := 0
		for i < len(books) {
			if books[i].Account.Username == username {
				sameAccountBooks = append(sameAccountBooks, books[i])
				books = append(books[:i], books[i+1:]...)
			} else {
				i++
			}
		}
		if len(sameAccountBooks) == 0 {
			break
		}
		l.Info(fmt.Sprintf("Download %d books for account %s", len(sameAccountBooks), sameAccountBooks[0].Account.Username))
		var downloadIdString strings.Builder
		for _, book := range sameAccountBooks {
			if downloadIdString.Len() > 0 {
				downloadIdString.WriteString(",")
			}
			downloadIdString.WriteString(book.DataId)
			downloadIdString.WriteString("=")
			downloadIdString.WriteString(filepath.Join(baseDir, book.UUID+".pvf"))
		}
		acc := sameAccountBooks[0].Account
		cmd := exec.Command("./integrations/d4s", "download", downloadIdString.String(), acc.Username, acc.Password)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		cmd.Start()
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				l.Info(scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				l.Err(fmt.Sprintf("Error occoured while downloading book: %s", scanner.Text()))
			}
		}()

		go func() {
			for {
				time.Sleep(180 * time.Second)
				err := rescanForDBInsert(baseDir, copyBooks)
				if err != nil {
					l.Err(fmt.Sprintf("Failed to rescan for books: %s", err.Error()))
				}
			}
		}()

		cmd.Wait()
		err := rescanForDBInsert(baseDir, copyBooks)
		if err != nil {
			l.Err(fmt.Sprintf("Failed to rescan for books: %s", err.Error()))
		}
	}
	return nil
}

func rescanForDBInsert(dir string, books []Book) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %v", dir, err)
	}
	for _, file := range files {
		for _, book := range books {
			if file.Name() == book.UUID+".pvf" {
				err := repo.Digi4SchoolBook.Save(&entity.Digi4SchoolBook{
					UUID:      book.UUID,
					BookName:  book.Name,
					BookID:    book.DataId,
					AccountID: book.Account.ID,
				})
				if err != nil {
					return fmt.Errorf("failed to save book %s: %v", book.UUID, err)
				}
			}
		}
	}

	return nil
}
func ListBooksForAccount(acc *entity.Digi4SchoolAccount) ([]Book, error) {
	cmd := exec.Command("./integrations/d4s", "list", acc.Username, acc.Password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to execute list for user %s: %v, output: %s", acc.Username, err, string(output))
		return nil, err
	}

	outStr := string(output)
	outStr = strings.TrimSpace(outStr)
	var books []Book
	if err := json.Unmarshal([]byte(outStr), &books); err != nil {
		log.Printf("failed to unmarshal list for user %s: %v, output: %s", acc.Username, err, string(output))
		return nil, err
	}
	for i, _ := range books {
		books[i].Account = acc
		books[i].UUID = uuid.NewString()
	}
	return books, nil
}
