package service

import (
	"context"
	"fmt"

	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/dto"
	searchRepo "github.com/YoungFlores/Case_Go/Profile/internal/search/repository"
)

type SearchServiceInterface interface {
	SearchProfileService(ctx context.Context, req dto.SearchDTO, helpers dto.SearchHelpersDTO) ([]models.Profile, error)
	SearchByFioService(ctx context.Context, req dto.SearchByFIODTO, helpers dto.SearchHelpersDTO) ([]models.Profile, error)
}

type SearchService struct {
	repo searchRepo.SearchRepo
}

func NewSearchService(repo searchRepo.SearchRepo) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

func (s *SearchService) SearchProfileService(ctx context.Context, req dto.SearchDTO, helpers dto.SearchHelpersDTO) ([]models.Profile, error) {
	var limit, page uint64 = 10, 1
	if helpers.Limit != nil {
		limit = *helpers.Limit
	}
	if helpers.Page != nil {
		page = *helpers.Page
	}

	offset := (page - 1) * limit

	if req.MinAge == nil && req.MaxAge != nil {
		age18 := 18
		req.MinAge = &age18
	}

	if req.MaxAge == nil && req.MinAge != nil {
		age100 := 100
		req.MaxAge = &age100
	}

	categories := []string{"id", "age", "user_id", "sex", "city"}
	orderBy := "id"
	if helpers.OrderBy != nil {
		for _, category := range categories {
			if *helpers.OrderBy == category {
				orderBy = category
				break
			}
		}
	}

	orderDirection := "ASC"
	if helpers.OrderDirection != nil && (*helpers.OrderDirection == "ASC" || *helpers.OrderDirection == "DESC") {
		orderDirection = *helpers.OrderDirection
	}

	users, err := s.repo.SearchProfile(ctx, req, limit, offset, orderBy, orderDirection)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (s *SearchService) SearchByFioService(ctx context.Context, req dto.SearchByFIODTO, helpers dto.SearchHelpersDTO) ([]models.Profile, error) {
	var limit uint64 = 10
	if helpers.Limit != nil && *helpers.Limit > 0 {
		limit = *helpers.Limit
	}

	var page uint64 = 1
	if helpers.Page != nil && *helpers.Page > 0 {
		page = *helpers.Page
	}

	offset := (page - 1) * limit

	if !isFioPresent(req) {
		return nil, fmt.Errorf("search criteria (name, surname or patronymic) must be provided")
	}

	profiles, err := s.repo.SearchByFio(ctx, req, limit, offset)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func isFioPresent(req dto.SearchByFIODTO) bool {
	if req.Name != nil && *req.Name != "" {
		return true
	}
	if req.Surname != nil && *req.Surname != "" {
		return true
	}
	if req.Patronymic != nil && *req.Patronymic != "" {
		return true
	}
	return false
}
