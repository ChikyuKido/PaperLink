package repo

import (
	"gorm.io/gorm"
	"paperlink/db"
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

	result := r.db.Where("id = ?", id).First(&entity)
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
