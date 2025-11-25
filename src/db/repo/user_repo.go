package repo

import (
	"paperlink/db/entity"
)

type UserRepo struct {
	*Repository[entity.User]
}

func newUserRepo() *UserRepo {
	return &UserRepo{NewRepository[entity.User]()}
}

var User = newUserRepo()

func (n *UserRepo) GetUserByName(name string) ([]entity.User, error) {
	var users []entity.User
	err := n.db.Where("Name = ?", name).Find(&users).Error
	return users, err
}

func (n *DocumentRepo) GetOwnedDocuments(userId int) ([]entity.Document, error) {
	var documents []entity.Document
	err := n.db.Where("OwnerID = ?", userId).Find(&documents).Error
	return documents, err
}
