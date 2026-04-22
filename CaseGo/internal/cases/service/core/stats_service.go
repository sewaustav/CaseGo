package service

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
)

func (s *CaseGoCoreService) GetStatsService(ctx context.Context) (*dto.StatsResponse, error) {
	totalCases, err := s.caseGoRepo.CountCases(ctx)
	if err != nil {
		return nil, err
	}

	totalDialogs, err := s.dialogRepo.CountDialogs(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.StatsResponse{
		TotalCases:   totalCases,
		TotalDialogs: totalDialogs,
	}, nil
}
