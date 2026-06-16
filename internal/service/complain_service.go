package service

import (
	"absensi-sppg/internal/repository"
	"math"
)

type ComplainService struct {
	repo repository.ComplainRepository
}

func NewComplainService(repo repository.ComplainRepository) ComplainService {
	return ComplainService{
		repo: repo,
	}
}

func (s *ComplainService) GetComplains(page int, status string, types string, parameter_flat string, parameter_zero string, ispu string, mutu_air string, validate string, accuweather_values string, per_page int, all_error string) (map[string]interface{}, error) {
	limit := 20

	data, total, err := s.repo.GetAll(page, limit, status, types, parameter_flat, parameter_zero, ispu, mutu_air, validate, accuweather_values, per_page, all_error)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data":       data,
		"page":       page,
		"total":      total,
		"total_page": totalPage,
	}

	return response, nil
}
