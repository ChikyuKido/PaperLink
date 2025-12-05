package repo

import (
	"paperlink/db/entity"
)

type DocumentRepo struct {
	*Repository[entity.Document]
}

func newDocumentRepo() *DocumentRepo {
	return &DocumentRepo{NewRepository[entity.Document]()}
}

var DocumentAction = newDocumentRepo()

func (n *DocumentRepo) GetAnnotationsById(documentID int) ([]entity.Annotation, error) {
	var annotations []entity.Annotation
	err := n.db.Where("ID = ?", documentID).Find(&annotations).Error
	return annotations, err
}
