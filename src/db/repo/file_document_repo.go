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

func (f *FileDocumentRepo) GetByUUID(uuid string) *entity.FileDocument {
	var doc entity.FileDocument
	tx := f.db.Where("uuid = ?", uuid).First(&doc)
	if tx.Error != nil {
		return nil
	}
	return &doc
}
