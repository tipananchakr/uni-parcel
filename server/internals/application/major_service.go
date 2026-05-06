package application

import (
	"context"
	"time"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"github.com/tipananchakr/uni-parcel/internals/core/ports"
)

type MajorService struct {
	majorRepository ports.MajorRepository
}

func NewMajorService(majoRepository ports.MajorRepository) *MajorService {
	return &MajorService{
		majorRepository: majoRepository,
	}
}

func (m *MajorService) GetAllMajors(ctx context.Context) ([]*domain.Major, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.majorRepository.GetAll(ctx)
}

func (m *MajorService) GetMajorByID(ctx context.Context, id string) (*domain.Major, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.majorRepository.GetByID(ctx, id)
}

func (m *MajorService) CreateMajor(ctx context.Context, major *domain.Major) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.majorRepository.Create(ctx, major)
}

func (m *MajorService) UpdateMajor(ctx context.Context, id string, update domain.MajorUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.majorRepository.Update(ctx, id, update)
}

func (m *MajorService) DeleteMajor(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.majorRepository.Delete(ctx, id)
}
