package auto

import (
	"context"
	"log"

	"github.com/bochkov/gobot/internal/db"
)

type Code struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type Region struct {
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Codes []string `json:"codes"`
}

func findCodesByRegion(id int) []string {
	ctx := context.Background()
	query := `select c.val from code c where c.region_id = $1`
	rows, err := db.GetPool().Query(ctx, query, id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	codes := make([]string, 0)
	for rows.Next() {
		var c string
		if err = rows.Scan(&c); err == nil {
			codes = append(codes, c)
		}
	}
	return codes
}

func FindRegionByName(name string) ([]Region, error) {
	ctx := context.Background()
	log.Printf("search region by name = '%s'", name)
	query := `select r.id, r.name from region r where r.name ILIKE '%' || $1 || '%'`
	rows, err := db.GetPool().Query(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	regions := make([]Region, 0)
	for rows.Next() {
		var r Region
		if err = rows.Scan(&r.Id, &r.Name); err == nil {
			r.Codes = findCodesByRegion(r.Id)
			regions = append(regions, r)
		}
	}
	return regions, nil
}

func FindRegionByCode(code string) (*Region, error) {
	ctx := context.Background()
	log.Printf("search region by code = '%s'", code)
	query := `select r.id rid, r.name rname from region r, code c where r.id = c.region_id and c.val = $1`
	var r Region
	if err := db.GetPool().QueryRow(ctx, query, code).Scan(&r.Id, &r.Name); err != nil {
		return nil, err
	}
	r.Codes = findCodesByRegion(r.Id)
	return &r, nil
}
