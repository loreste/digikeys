package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

type EnrollmentService struct {
	enrollmentRepo ports.EnrollmentRepository
	citizenRepo    ports.CitizenRepository
}

func NewEnrollmentService(enrollmentRepo ports.EnrollmentRepository, citizenRepo ports.CitizenRepository) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		citizenRepo:    citizenRepo,
	}
}

func (s *EnrollmentService) CreateEnrollment(ctx context.Context, enrollment *domain.Enrollment) error {
	enrollment.ID = uuid.New().String()
	enrollment.SyncStatus = "synced"
	enrollment.ReviewStatus = "pending"
	enrollment.EnrolledAt = time.Now()
	enrollment.CreatedAt = time.Now()
	enrollment.UpdatedAt = time.Now()

	return s.enrollmentRepo.Create(ctx, enrollment)
}

type MobileSyncRequest struct {
	Enrollments []domain.Enrollment `json:"enrollments"`
}

type MobileSyncResponse struct {
	Synced int      `json:"synced"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors,omitempty"`
}

func (s *EnrollmentService) SyncFromMobile(ctx context.Context, batch []domain.Enrollment) (*MobileSyncResponse, error) {
	resp := &MobileSyncResponse{}

	for i := range batch {
		e := &batch[i]
		if e.ID == "" {
			e.ID = uuid.New().String()
		}
		e.SyncStatus = "synced"
		now := time.Now()
		e.SyncedAt = &now
		e.ReviewStatus = "pending"
		e.CreatedAt = now
		e.UpdatedAt = now

		if err := s.enrollmentRepo.Create(ctx, e); err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, fmt.Sprintf("enrollment %s: %v", e.ID, err))
			continue
		}
		resp.Synced++
	}

	if resp.Failed > 0 && resp.Synced == 0 {
		return resp, domain.ErrSyncFailed
	}

	return resp, nil
}

func (s *EnrollmentService) ReviewEnrollment(ctx context.Context, id, status, reviewedBy, notes string) error {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if enrollment.ReviewStatus != "pending" && enrollment.ReviewStatus != "needs_correction" {
		return fmt.Errorf("%w: enrollment already reviewed", domain.ErrInvalidInput)
	}

	if status != "approved" && status != "rejected" && status != "needs_correction" {
		return fmt.Errorf("%w: invalid review status", domain.ErrInvalidInput)
	}

	now := time.Now()
	return s.enrollmentRepo.UpdateReviewStatus(ctx, id, status, reviewedBy, notes, &now)
}

func (s *EnrollmentService) ListByAgent(ctx context.Context, agentID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return s.enrollmentRepo.ListByAgent(ctx, agentID, page, pageSize)
}

func (s *EnrollmentService) ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return s.enrollmentRepo.ListByEmbassy(ctx, embassyID, page, pageSize)
}

func (s *EnrollmentService) List(ctx context.Context, filter domain.EnrollmentFilter) ([]*domain.Enrollment, int, error) {
	return s.enrollmentRepo.List(ctx, filter)
}

func (s *EnrollmentService) GetByID(ctx context.Context, id string) (*domain.Enrollment, error) {
	return s.enrollmentRepo.GetByID(ctx, id)
}
