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

func (n *UserRepo) DoesUserByNameExist(name string) bool {
	var count int64
	err := n.db.Model(&entity.User{}).Where("Username = ?", name).Count(&count).Error
	if err != nil {
		return true
	}
	return count > 0
}

func (n *UserRepo) GetUserByName(name string) (*entity.User, error) {
	var user entity.User
	err := n.db.Where("Username = ?", name).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (n *DocumentRepo) GetOwnedDocuments(userId int) ([]entity.Document, error) {
	var documents []entity.Document
	err := n.db.Where("OwnerID = ?", userId).Find(&documents).Error
	return documents, err
}
