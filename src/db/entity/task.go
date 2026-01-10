package entity

type TaskStatus string

const (
	RUNNING   TaskStatus = "RUNNING"
	FAILED    TaskStatus = "FAILED"
	COMPLETED TaskStatus = "COMPLETED"
)

type Task struct {
	ID        string     `gorm:"primary_key" json:"id"`
	Status    TaskStatus `json:"status"`
	Name      string     `gorm:"not null" json:"name"`
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
}
