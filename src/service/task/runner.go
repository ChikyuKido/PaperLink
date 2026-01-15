package task

import (
	"fmt"
	"os"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"path/filepath"
	"strings"
	"sync"
)

var (
	taskStore   = make(map[string]*TaskRunner)
	taskStoreMu sync.RWMutex
	dataDir     = "./data/tasks"
)

type TaskRunner struct {
	Task     *entity.Task
	logs     []string
	logMu    sync.Mutex
	Complete func() error
	Fail     func() error
}

type TaskInfo struct {
	UUID   string            `json:"uuid"`
	Name   string            `json:"name"`
	Status entity.TaskStatus `json:"status"`
	Lines  []string          `json:"lines,omitempty"`
}

func Init() {
	_ = os.MkdirAll(dataDir, os.ModePerm)
}

func CreateNewTask(name string) (*TaskRunner, error) {
	task, err := repo.Task.StartTask(name)
	if err != nil {
		return nil, err
	}

	runner := &TaskRunner{
		Task: task,
		logs: []string{},
	}
	runner.Complete = func() error { return completeTask(runner) }
	runner.Fail = func() error { return failTask(runner) }

	taskStoreMu.Lock()
	taskStore[task.ID] = runner
	taskStoreMu.Unlock()

	return runner, nil
}

func (tr *TaskRunner) log(level, msg string) {
	tr.logMu.Lock()
	defer tr.logMu.Unlock()
	line := fmt.Sprintf("[%s] %s", level, msg)
	tr.logs = append(tr.logs, line)
}

func (tr *TaskRunner) Info(msg string)     { tr.log("INFO", msg) }
func (tr *TaskRunner) Warn(msg string)     { tr.log("WARN", msg) }
func (tr *TaskRunner) Err(msg string)      { tr.log("ERROR", msg) }
func (tr *TaskRunner) Critical(msg string) { tr.log("CRITICAL", msg) }
func (tr *TaskRunner) ReplaceLastInfo(msg string) {
	if len(tr.logs) == 0 {
		return
	}
	if strings.Contains(tr.logs[len(tr.logs)-1], "INFO") {
		tr.logs = tr.logs[:len(tr.logs)-1]
		tr.Info(msg)
	}
}

func completeTask(tr *TaskRunner) error {
	tr.logMu.Lock()
	lines := append([]string(nil), tr.logs...)
	tr.logMu.Unlock()

	filePath := filepath.Join(dataDir, tr.Task.ID+".log")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, _ = f.WriteString(line + "\n")
	}

	return repo.Task.FinishTask(tr.Task)
}

func failTask(tr *TaskRunner) error {
	tr.logMu.Lock()
	lines := append([]string(nil), tr.logs...)
	tr.logMu.Unlock()

	filePath := filepath.Join(dataDir, tr.Task.ID+".log")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, _ = f.WriteString(line + "\n")
	}

	return repo.Task.FailTask(tr.Task)
}

func GetTaskInfo(uuid string) (*TaskInfo, error) {
	taskStoreMu.RLock()
	runner, ok := taskStore[uuid]
	taskStoreMu.RUnlock()

	var lines []string
	var task *entity.Task

	if ok {
		runner.logMu.Lock()
		lines = append([]string(nil), runner.logs...)
		runner.logMu.Unlock()
		task = runner.Task
	} else {
		filePath := filepath.Join(dataDir, uuid+".log")
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		lines = append([]string(nil), string(content))
		task, err = repo.Task.Get(uuid)
		if err != nil {
			return nil, err
		}
	}

	return &TaskInfo{
		UUID:   uuid,
		Name:   task.Name,
		Status: task.Status,
		Lines:  lines,
	}, nil
}
func ListTasks() ([]*TaskInfo, error) {
	taskStoreMu.RLock()
	list := make([]*TaskInfo, 0, len(taskStore))
	for _, runner := range taskStore {
		list = append(list, &TaskInfo{
			UUID:   runner.Task.ID,
			Name:   runner.Task.Name,
			Status: runner.Task.Status,
		})
	}
	taskStoreMu.RUnlock()

	completedTasks, err := repo.Task.ListCompletedOrFailed()
	if err != nil {
		return nil, err
	}
	for _, t := range completedTasks {
		list = append(list, &TaskInfo{
			UUID:   t.ID,
			Name:   t.Name,
			Status: t.Status,
		})
	}

	return list, nil
}
