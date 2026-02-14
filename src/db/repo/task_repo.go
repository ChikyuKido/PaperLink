package repo

import (
	"github.com/google/uuid"
	"paperlink/db/entity"
	"time"
)

type TaskRepo struct {
	*Repository[entity.Task]
}

func newTaskRepo() *TaskRepo {
	return &TaskRepo{NewRepository[entity.Task]()}
}

var Task = newTaskRepo()

func (repo *TaskRepo) StartTask(name string) (*entity.Task, error) {
	task := entity.Task{
		ID:        uuid.New().String(),
		Status:    entity.RUNNING,
		Name:      name,
		StartTime: time.Now().Unix(),
	}
	err := repo.Save(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TaskRepo) FinishTask(task *entity.Task) error {
	task.Status = entity.COMPLETED
	task.EndTime = time.Now().Unix()
	err := repo.Save(task)
	if err != nil {
		return err
	}
	return nil
}
func (repo *TaskRepo) FailTask(task *entity.Task) error {
	task.Status = entity.FAILED
	task.EndTime = time.Now().Unix()
	err := repo.Save(task)
	if err != nil {
		return err
	}
	return nil
}
func (repo *TaskRepo) StopTask(task *entity.Task) error {
	task.Status = entity.STOPPED
	task.EndTime = time.Now().Unix()
	err := repo.Save(task)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepo) ListCompletedOrFailed() ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := repo.db.Where("status IN ?", []entity.TaskStatus{entity.COMPLETED, entity.FAILED, entity.STOPPED}).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
