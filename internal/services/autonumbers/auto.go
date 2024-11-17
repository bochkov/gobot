package autonumbers

import (
	"context"

	"github.com/bochkov/gobot/internal/services"
)

type Code struct {
	Id    int    `json:"id" db:"id"`
	Value string `json:"value" db:"val"`
}

type Region struct {
	Id    int      `json:"id" db:"id"`
	Name  string   `json:"name" db:"name"`
	Codes []string `json:"codes"`
}

type Repository interface {
	FindRegionByCode(ctx context.Context, code string) (*Region, error)
	FindRegionByName(ctx context.Context, name string) ([]Region, error)
}

type Service interface {
	services.Service
	FindRegionByCode(ctx context.Context, code string) (*Region, error)
	FindRegionByName(ctx context.Context, name string) ([]Region, error)
}
