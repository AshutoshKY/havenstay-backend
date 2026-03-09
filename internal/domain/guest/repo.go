package guest

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, guestID int64) (*Entity, error)
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

func (r *MySQLRepository) GetByID(ctx context.Context, guestID int64) (*Entity, error) {
	var model Model
	if err := r.db.WithContext(ctx).First(&model, guestID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}
	return ModelToEntity(&model), nil
}
