package application

import (
	"context"

	"github.com/digikeys/backend/internal/ports"
)

type StatisticsService struct {
	querier ports.StatisticsQuerier
}

func NewStatisticsService(querier ports.StatisticsQuerier) *StatisticsService {
	return &StatisticsService{querier: querier}
}

type DashboardStats struct {
	TotalCitizens       int            `json:"totalCitizens"`
	CardsByStatus       map[string]int `json:"cardsByStatus"`
	EnrollmentsByStatus map[string]int `json:"enrollmentsByStatus"`
	TransfersTotal      int64          `json:"transfersTotal"`
	FSBTotal            int64          `json:"fsbTotal"`
	CitizensByCountry   map[string]int `json:"citizensByCountry,omitempty"`
}

func (s *StatisticsService) GetDashboard(ctx context.Context, embassyID string) (*DashboardStats, error) {
	stats := &DashboardStats{}

	totalCitizens, err := s.querier.TotalCitizens(ctx, embassyID)
	if err != nil {
		return nil, err
	}
	stats.TotalCitizens = totalCitizens

	cardsByStatus, err := s.querier.CardsByStatus(ctx, embassyID)
	if err != nil {
		return nil, err
	}
	stats.CardsByStatus = cardsByStatus

	enrollmentsByStatus, err := s.querier.EnrollmentsByStatus(ctx, embassyID)
	if err != nil {
		return nil, err
	}
	stats.EnrollmentsByStatus = enrollmentsByStatus

	transfersTotal, err := s.querier.TransfersTotal(ctx, embassyID)
	if err != nil {
		return nil, err
	}
	stats.TransfersTotal = transfersTotal

	fsbTotal, err := s.querier.FSBTotal(ctx, embassyID)
	if err != nil {
		return nil, err
	}
	stats.FSBTotal = fsbTotal

	if embassyID == "" {
		citizensByCountry, err := s.querier.CitizensByCountry(ctx)
		if err != nil {
			return nil, err
		}
		stats.CitizensByCountry = citizensByCountry
	}

	return stats, nil
}
