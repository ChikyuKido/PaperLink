package repo

import (
	"paperlink/db/entity"
)

type TagRepo struct {
	*Repository[entity.Tag]
}

func newTagRepo() *TagRepo {
	return &TagRepo{NewRepository[entity.Tag]()}
}

var Tag = newTagRepo()

func (n *DocumentRepo) GetDocumentsWithTag(tagId int) ([]entity.Document, error) {
	var documents []entity.Document
	err := n.db.Where("ID = ?", tagId).Find(&documents).Error
	return documents, err
}
