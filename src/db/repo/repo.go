package repo

import (
	"paperlink/db"

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

func (r *Repository[T]) Get(id int) (*T, error) {
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

func (r *Repository[T]) Delete(id int) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}
