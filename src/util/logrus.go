package util

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type GroupFormatter struct{}

func (f *GroupFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	group, ok := entry.Data["group"].(string)
	if !ok {
		group = "General"
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := fmt.Sprintf("[%s] [%s] [%s] %s\n", timestamp, group, level, entry.Message)
	return []byte(msg), nil
}

func GroupLog(group string) *logrus.Entry {
	return logrus.WithField("group", group)
}
