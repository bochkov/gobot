package autonumbers

import (
	"context"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	return &service{
		repo,
		time.Duration(3) * time.Second,
	}
}

func (s *service) Description() string {
	return "Автомобильные коды регионов РФ"
}

func (s *service) FindRegionByCode(c context.Context, code string) (*Region, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	return s.Repository.FindRegionByCode(ctx, code)
}

func (s *service) FindRegionByName(c context.Context, name string) ([]Region, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	return s.Repository.FindRegionByName(ctx, name)
}
