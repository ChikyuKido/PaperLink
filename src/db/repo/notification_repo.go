package repo

import (
	"paperlink/db/entity"
)

type NotificationRepo struct {
	*Repository[entity.Notification]
}

func newNotificationRepo() *NotificationRepo {
	return &NotificationRepo{NewRepository[entity.Notification]()}
}

var Notification = newNotificationRepo()

func (n *NotificationRepo) GetNotificationsForUser(userID int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := n.db.Where("ID = ?", userID).Find(&notifications).Error
	return notifications, err
}
