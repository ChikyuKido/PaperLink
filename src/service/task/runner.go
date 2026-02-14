package task

import (
	"errors"
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

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskNotRunning = errors.New("task is not running")
)

type TaskRunner struct {
	Task        *entity.Task
	logs        []string
	logMu       sync.Mutex
	stopHandler func(*TaskRunner) error
	Complete    func() error
	Fail        func() error
	Stop        func() error
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

func CreateNewTask(name string, stopHandler ...func(*TaskRunner) error) (*TaskRunner, error) {
	task, err := repo.Task.StartTask(name)
	if err != nil {
		return nil, err
	}

	var stopper func(*TaskRunner) error
	if len(stopHandler) > 0 {
		stopper = stopHandler[0]
	}

	runner := &TaskRunner{
		Task:        task,
		logs:        []string{},
		stopHandler: stopper,
	}
	runner.Complete = func() error { return completeTask(runner) }
	runner.Fail = func() error { return failTask(runner) }
	runner.Stop = func() error { return stopTask(runner) }

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
	if tr.Task.Status != entity.RUNNING {
		return nil
	}
	tr.logMu.Lock()
	lines := append([]string(nil), tr.logs...)
	tr.logMu.Unlock()

	if err := writeLogFile(tr.Task.ID, lines); err != nil {
		return err
	}
	if err := repo.Task.FinishTask(tr.Task); err != nil {
		return err
	}
	removeTask(tr.Task.ID)
	return nil
}

func failTask(tr *TaskRunner) error {
	if tr.Task.Status != entity.RUNNING {
		return nil
	}
	tr.logMu.Lock()
	lines := append([]string(nil), tr.logs...)
	tr.logMu.Unlock()

	if err := writeLogFile(tr.Task.ID, lines); err != nil {
		return err
	}

	if err := repo.Task.FailTask(tr.Task); err != nil {
		return err
	}
	removeTask(tr.Task.ID)
	return nil
}

func stopTask(tr *TaskRunner) error {
	if tr.Task.Status != entity.RUNNING {
		return ErrTaskNotRunning
	}

	if tr.stopHandler != nil {
		if err := tr.stopHandler(tr); err != nil {
			tr.Err(fmt.Sprintf("stop handler failed: %v", err))
		}
	}

	tr.logMu.Lock()
	lines := append([]string(nil), tr.logs...)
	tr.logMu.Unlock()

	if err := writeLogFile(tr.Task.ID, lines); err != nil {
		return err
	}
	if err := repo.Task.StopTask(tr.Task); err != nil {
		return err
	}
	removeTask(tr.Task.ID)
	return nil
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
		raw := strings.TrimSpace(string(content))
		if raw != "" {
			lines = strings.Split(raw, "\n")
		} else {
			lines = []string{}
		}
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

func StopTask(uuid string) error {
	taskStoreMu.RLock()
	runner, ok := taskStore[uuid]
	taskStoreMu.RUnlock()
	if !ok {
		if t, err := repo.Task.Get(uuid); err != nil || t == nil {
			return ErrTaskNotFound
		}
		return ErrTaskNotRunning
	}
	return runner.Stop()
}

func (tr *TaskRunner) IsRunning() bool {
	return tr.Task.Status == entity.RUNNING
}

func writeLogFile(taskID string, lines []string) error {
	filePath := filepath.Join(dataDir, taskID+".log")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, _ = f.WriteString(line + "\n")
	}
	return nil
}

func removeTask(taskID string) {
	taskStoreMu.Lock()
	delete(taskStore, taskID)
	taskStoreMu.Unlock()
}
