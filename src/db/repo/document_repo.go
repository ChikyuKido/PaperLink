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

var Document = newDocumentRepo()

func (n *DocumentRepo) GetAnnotationsById(documentID int) ([]entity.Annotation, error) {
	var annotations []entity.Annotation
	err := n.db.Where("ID = ?", documentID).Find(&annotations).Error
	return annotations, err
}

func (r *DocumentRepo) GetAllByUserId(userId int) ([]entity.Document, error) {
	var result []entity.Document
	err := r.db.Where("user_id = ?", userId).Find(&result).Error
	return result, err
}
