package property

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository interface for the Property domain
type Repository interface {
	Create(ctx context.Context, entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, id int64) (*Entity, error)
	ListByHostID(ctx context.Context, hostID int64, locationFilter string) ([]*Entity, error)
	ListAll(ctx context.Context) ([]*Entity, error)
}

// MySQLRepository is the GORM implementation of the Property Repository
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
			return nil, errors.New("property not found")
		}
		return nil, err
	}
	return ModelToEntity(&model), nil
}

func (r *MySQLRepository) ListByHostID(ctx context.Context, hostID int64, locationFilter string) ([]*Entity, error) {
	var models []Model
	query := r.db.WithContext(ctx).Where("host_id = ?", hostID)
	if locationFilter != "" {
		query = query.Where("location LIKE ?", "%"+locationFilter+"%")
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	var entities []*Entity
	for _, m := range models {
		mCopy := m // fix loop var pointer aliasing
		entities = append(entities, ModelToEntity(&mCopy))
	}
	return entities, nil
}

func (r *MySQLRepository) ListAll(ctx context.Context) ([]*Entity, error) {
	var models []Model
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	var entities []*Entity
	for _, m := range models {
		mCopy := m
		entities = append(entities, ModelToEntity(&mCopy))
	}
	return entities, nil
}
