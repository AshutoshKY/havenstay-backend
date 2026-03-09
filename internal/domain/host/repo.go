package host

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository interface for the Host domain
type Repository interface {
	Create(ctx context.Context, entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, id int64) (*Entity, error)
}

// MySQLRepository is the GORM implementation of the Host Repository
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository creates a new MySQLRepository
func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Create inserts a new host into the database
func (r *MySQLRepository) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	model := EntityToModel(entity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return ModelToEntity(model), nil
}

// GetByID retrieves a host by its ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Entity, error) {
	var model Model
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("host not found")
		}
		return nil, err
	}
	return ModelToEntity(&model), nil
}
