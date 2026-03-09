package wishlist

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Add(ctx context.Context, entity *Entity) (*Entity, error)
	ListByGuestID(ctx context.Context, guestID int64) ([]*Entity, error)
}

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Add(ctx context.Context, entity *Entity) (*Entity, error) {
	model := EntityToModel(entity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return ModelToEntity(model), nil
}

func (r *MySQLRepository) ListByGuestID(ctx context.Context, guestID int64) ([]*Entity, error) {
	var models []Model
	if err := r.db.WithContext(ctx).Where("guest_id = ?", guestID).Find(&models).Error; err != nil {
		return nil, err
	}

	var entities []*Entity
	for _, m := range models {
		mCopy := m
		entities = append(entities, ModelToEntity(&mCopy))
	}
	return entities, nil
}
