package repo

import (
	"paperlink/db"
	"paperlink/db/entity"
	"strings"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{db: db.DB()}
}

func (r *Repository[T]) Save(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *Repository[T]) SaveList(entities []*T) error {
	return r.db.Save(&entities).Error
}

func (r *Repository[T]) Get(id any) (*T, error) {
	var entity T
	result := r.db.First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (r *Repository[T]) GetList() ([]T, error) {
	var entities []T
	result := r.db.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

func (r *Repository[T]) Delete(id any) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}
func (r *Repository[T]) GetByIDs(ids []any) ([]T, error) {
	var entities []T
	result := r.db.Where("id IN ?", ids).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
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

// Filter: tags = "alle tags müssen vorhanden sein" (AND)
// search matched in name OR description
// dirID:
//   - nil  => keine directory filter (alles)
//   - &0   => root (DirectoryID IS NULL)
//   - &x   => exact directory (DirectoryID = x)
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
		// AND-Filter: doc muss alle tags enthalten
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

// Save mit Many2Many korrekt (Tags updaten)
func (r *DocumentRepo) SaveWithAssociations(doc *entity.Document) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(doc).Error
}
