package entity

type TaskStatus string

const (
	RUNNING   TaskStatus = "RUNNING"
	FAILED    TaskStatus = "FAILED"
	COMPLETED TaskStatus = "COMPLETED"
)

type Task struct {
	ID int `gorm:"primary_key" json:"id"`
}
