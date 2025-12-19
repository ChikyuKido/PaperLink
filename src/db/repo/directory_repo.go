package repo

import (
	"paperlink/db/entity"
)

type DirectoryRepo struct {
	*Repository[entity.Directory]
}

func newDirectoryRepo() *DirectoryRepo {
	return &DirectoryRepo{NewRepository[entity.Directory]()}
}

var Directory = newDirectoryRepo()

func (r *DirectoryRepo) GetAllByUserId(userId int) ([]entity.Directory, error) {
	var result []entity.Directory
	err := r.db.Where("user_id = ?", userId).Find(&result).Error
	return result, err
}
