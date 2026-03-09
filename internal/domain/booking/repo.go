package booking

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, id int64) (*Entity, error)
	ListByGuestID(ctx context.Context, guestID int64) ([]*Entity, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	model := EntityToModel(entity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return ModelToEntity(model), nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Entity, error) {
	var model Model
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return ModelToEntity(&model), nil
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

func (r *MySQLRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&Model{}).Where("booking_id = ?", id).Update("status", status).Error
}
