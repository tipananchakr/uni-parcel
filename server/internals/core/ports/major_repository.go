package ports

import (
	"context"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
)

type MajorRepository interface {
	GetAll(ctx context.Context) ([]*domain.Major, error)
	GetByID(ctx context.Context, id string) (*domain.Major, error)
	Create(ctx context.Context, major *domain.Major) error
	Update(ctx context.Context, id string, update domain.MajorUpdate) error
	Delete(ctx context.Context, id string) error
}
