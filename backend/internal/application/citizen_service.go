package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

type CitizenService struct {
	citizenRepo ports.CitizenRepository
}

func NewCitizenService(citizenRepo ports.CitizenRepository) *CitizenService {
	return &CitizenService{citizenRepo: citizenRepo}
}

func (s *CitizenService) Create(ctx context.Context, citizen *domain.Citizen) error {
	if citizen.NationalID != "" {
		existing, _ := s.citizenRepo.GetByNationalID(ctx, citizen.NationalID)
		if existing != nil {
			return fmt.Errorf("%w: citizen with this national ID already exists", domain.ErrAlreadyExists)
		}
	}

	citizen.ID = uuid.New().String()
	citizen.Status = "active"
	citizen.CreatedAt = time.Now()
	citizen.UpdatedAt = time.Now()

	return s.citizenRepo.Create(ctx, citizen)
}

func (s *CitizenService) GetByID(ctx context.Context, id string) (*domain.Citizen, error) {
	return s.citizenRepo.GetByID(ctx, id)
}

func (s *CitizenService) Search(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return s.citizenRepo.Search(ctx, filter)
}

func (s *CitizenService) Update(ctx context.Context, citizen *domain.Citizen) error {
	existing, err := s.citizenRepo.GetByID(ctx, citizen.ID)
	if err != nil {
		return err
	}

	citizen.CreatedAt = existing.CreatedAt
	citizen.UpdatedAt = time.Now()

	return s.citizenRepo.Update(ctx, citizen)
}

func (s *CitizenService) List(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return s.citizenRepo.List(ctx, filter)
}
