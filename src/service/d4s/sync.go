package d4s

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/service/task"
	"strings"
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
	for _, book := range accountBooks {
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
	l.Info(fmt.Sprintf("Found %d needed books", len(neededBooks)))

	err = l.Complete()
	if err != nil {
		log.Error("Failed to complete the sync task")
	}
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
	return books, nil
}
