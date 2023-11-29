package cbr

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

const (
	insCurrencyItemQuery     = `insert into currency_item (id, name, eng_name, nominal, parent_code, iso_num_code, iso_char_code) values ($1, $2, $3, $4, $5, $6, $7)`
	insCurrencyRateQuery     = `insert into currency_rate (date, name) VALUES ($1, $2) returning id`
	insCurrencyRateItemQuery = `insert into currency_rate_record (curr_id, rate_id, num_code, char_code, nominal, name, rate_value) values ($1, $2, $3, $4, $5, $6, $7)`
)

func (r *repository) FindRateRecordsByRateId(ctx context.Context, rateId int) ([]RateItem, error) {
	query :=
		`SELECT crr.id, crr.curr_id, crr.num_code, crr.char_code, crr.nominal, crr.name, crr.rate_value
		 FROM currency_rate_record crr 
		 WHERE crr.rate_id = $1`
	rows, err := r.db.Query(ctx, query, rateId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]RateItem, 0)
	for rows.Next() {
		var ri RateItem
		if err := rows.Scan(&ri.Id, &ri.CurID, &ri.NumCode, &ri.CharCode, &ri.Nominal, &ri.Name, &ri.Value); err != nil {
			log.Fatal(err)
		}
		items = append(items, ri)
	}
	return items, nil
}

func (r *repository) LatestRate(ctx context.Context) (*CurrRate, error) {
	query :=
		`SELECT cr.id, cr.date, cr.fetch_time, cr.name 
		 FROM currency_rate cr 
		 WHERE cr.date = (SELECT max(cr1.date) FROM currency_rate cr1)`
	row := r.db.QueryRow(ctx, query)

	var cr CurrRate
	if err := row.Scan(&cr.Id, &cr.Date, &cr.FetchDate, &cr.Name); err != nil {
		log.Fatal(err)
		return nil, err
	}

	records, err := r.FindRateRecordsByRateId(ctx, int(cr.Id))
	if err != nil {
		return nil, err
	}
	cr.RateItems = records
	return &cr, nil
}
