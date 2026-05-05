package application

import (
	"context"
	"time"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"github.com/tipananchakr/uni-parcel/internals/core/ports"
)

type DormService struct {
	dormRepository ports.DormRepository
}

func NewDormService(ctx context.Context, dormRepository ports.DormRepository) *DormService {
	return &DormService{
		dormRepository: dormRepository,
	}
}

func (d *DormService) GetAllDorms(ctx context.Context) ([]*domain.Dorm, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.dormRepository.GetAll(ctx)
}

func (d *DormService) GetDormByID(ctx context.Context, id string) (*domain.Dorm, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.dormRepository.GetByID(ctx, id)
}

func (d *DormService) CreateDorm(ctx context.Context, dorm *domain.Dorm) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.dormRepository.Create(ctx, dorm)
}

func (d *DormService) UpdateDorm(ctx context.Context, id string, update domain.DormUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.dormRepository.Update(ctx, id, update)
}

func (d *DormService) DeleteDorm(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.dormRepository.Delete(ctx, id)
}
