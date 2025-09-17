package postgres

import (
	"github.com/rauan06/realtime-map/analytics/internal/domain"
	"gorm.io/gorm"
)

type Repository[T domain.Entity] struct {
	db *gorm.DB
}

func New[T domain.Entity](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *Repository[T]) GetByID(id string, preload ...string) (*T, error) {
	var entity T
	query := r.db
	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *Repository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *Repository[T]) Delete(id string) error {
	var entity T
	return r.db.Delete(&entity, "id = ?", id).Error
}
