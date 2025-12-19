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
