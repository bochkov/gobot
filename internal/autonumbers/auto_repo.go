package autonumbers

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) findCodesByRegion(ctx context.Context, id int) ([]string, error) {
	query := `select val from code where region_id = $1`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	codes := make([]string, 0)
	for rows.Next() {
		var c string
		if err = rows.Scan(&c); err == nil {
			codes = append(codes, c)
		}
	}
	return codes, nil
}

func (r *repository) FindRegionByName(ctx context.Context, name string) ([]Region, error) {
	log.Printf("search region by name = '%s'", name)
	query := `select r.id, r.name from region r where r.name ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := make([]Region, 0)
	for rows.Next() {
		var region Region
		if err = rows.Scan(&region.Id, &region.Name); err == nil {
			codes, err := r.findCodesByRegion(ctx, region.Id)
			if err != nil {
				log.Fatalf("cannot get codes for region = %d", region.Id)
			} else {
				region.Codes = codes
			}
			regions = append(regions, region)
		}
	}
	return regions, nil
}

func (r *repository) FindRegionByCode(ctx context.Context, code string) (*Region, error) {
	log.Printf("search region by code = '%s'", code)
	query := `select r.id rid, r.name rname from region r, code cd where r.id = cd.region_id and cd.val = $1`
	row := r.db.QueryRow(ctx, query, code)

	var region Region
	if err := row.Scan(&region.Id, &region.Name); err != nil {
		return nil, err
	}
	codes, err := r.findCodesByRegion(ctx, region.Id)
	if err != nil {
		log.Fatalf("cannot get codes for region = %d", region.Id)
	}
	region.Codes = codes
	return &region, nil
}
