package repo

import (
	"paperlink/db/entity"
)

type FileDocumentRepo struct {
	*Repository[entity.FileDocument]
}

func newFileDocumentRepo() *FileDocumentRepo {
	return &FileDocumentRepo{NewRepository[entity.FileDocument]()}
}

var FileDocument = newFileDocumentRepo()
