// Package repository holds the GORM models and the repository implementations
// that satisfy the port interfaces declared in internal/service.
package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"go-starter-kit/internal/domain/example"
	examplesvc "go-starter-kit/internal/service/example"
	"go-starter-kit/internal/shared"
)

// ExampleModel is the GORM persistence shape for example.Example.
type ExampleModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (ExampleModel) TableName() string { return "examples" }

func (m *ExampleModel) toDomain() *example.Example {
	return &example.Example{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func exampleModelFromDomain(e *example.Example) *ExampleModel {
	return &ExampleModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

type exampleRepo struct {
	db *gorm.DB
}

// NewExampleRepository returns a GORM-backed example service Repository. Each
// method uses shared.GetDB, so a transaction placed on the context by TxManager
// is used automatically; otherwise it falls back to the default handle.
func NewExampleRepository(db *gorm.DB) examplesvc.Repository {
	return &exampleRepo{db: db}
}

func (r *exampleRepo) Create(ctx context.Context, e *example.Example) error {
	m := exampleModelFromDomain(e)
	if err := shared.GetDB(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *exampleRepo) FindByID(ctx context.Context, id string) (*example.Example, error) {
	var m ExampleModel
	err := shared.GetDB(ctx, r.db).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, example.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.toDomain(), nil
}

func (r *exampleRepo) List(ctx context.Context, offset, limit int) ([]example.Example, int64, error) {
	var total int64
	if err := shared.GetDB(ctx, r.db).Model(&ExampleModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var models []ExampleModel
	if err := shared.GetDB(ctx, r.db).Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]example.Example, len(models))
	for i := range models {
		out[i] = *models[i].toDomain()
	}
	return out, total, nil
}

func (r *exampleRepo) Update(ctx context.Context, e *example.Example) error {
	m := exampleModelFromDomain(e)
	if err := shared.GetDB(ctx, r.db).Save(m).Error; err != nil {
		return err
	}
	e.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *exampleRepo) Delete(ctx context.Context, id string) error {
	res := shared.GetDB(ctx, r.db).Delete(&ExampleModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return example.ErrNotFound
	}
	return nil
}
