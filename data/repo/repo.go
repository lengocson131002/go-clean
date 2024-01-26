package repo

import "gorm.io/gorm"

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *Repository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *Repository[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *Repository[T]) FindById(id any) (*T, error) {
	entity := new(T)
	err := r.db.Where("id = ?", id).Take(entity).Error
	return entity, err
}
