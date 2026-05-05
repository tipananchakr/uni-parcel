package ports

import (
	"context"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
)

type DormRepository interface {
	GetAll(ctx context.Context) ([]*domain.Dorm, error)
	GetByID(ctx context.Context, id string) (*domain.Dorm, error)
	Create(ctx context.Context, dorm *domain.Dorm) error
	Update(ctx context.Context, id string, update domain.DormUpdate) error
	Delete(ctx context.Context, id string) error
}
