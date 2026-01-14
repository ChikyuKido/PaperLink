package repo

import (
	"gorm.io/gorm"
	"paperlink/db/entity"
	"strings"
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

func (r *DocumentRepo) GetAllByUserIdWithFile(userId int) ([]entity.Document, error) {
	var result []entity.Document
	err := r.db.Preload("File").Where("user_id = ?", userId).Find(&result).Error
	return result, err
}

func (r *DocumentRepo) GetByUUIDWithFile(uuid string) *entity.Document {
	var doc entity.Document
	err := r.db.Preload("File").Where("uuid = ?", uuid).First(&doc).Error
	if err != nil {
		return nil
	}
	return &doc
}

func (r *DocumentRepo) GetByUUIDWithTagsAndFile(uuid string) *entity.Document {
	var doc entity.Document
	err := r.db.
		Preload("Tags").
		Preload("File").
		Where("uuid = ?", uuid).
		First(&doc).Error
	if err != nil {
		return nil
	}
	return &doc
}

func (r *DocumentRepo) Filter(userID int, tags []string, search string) ([]entity.Document, error) {
	q := r.db.
		Model(&entity.Document{}).
		Preload("Tags").
		Preload("File").
		Where("documents.user_id = ?", userID)

	if search != "" {
		like := "%" + strings.TrimSpace(search) + "%"
		q = q.Where("(documents.name LIKE ? OR documents.description LIKE ?)", like, like)
	}

	if len(tags) > 0 {

		q = q.
			Joins("JOIN document_tags dt ON dt.document_id = documents.id").
			Joins("JOIN tags t ON t.id = dt.tag_id").
			Where("t.name IN ?", tags).
			Group("documents.id").
			Having("COUNT(DISTINCT t.name) = ?", len(tags))
	}

	var docs []entity.Document
	err := q.Find(&docs).Error
	return docs, err
}

func (r *DocumentRepo) SaveWithAssociations(doc *entity.Document) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(doc).Error
}
