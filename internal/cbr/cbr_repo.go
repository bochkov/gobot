package cbr

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) FindRateRecordsByRateId(ctx context.Context, rateId int) ([]RateItem, error) {
	query :=
		`SELECT crr.id, crr.curr_id, crr.num_code, crr.char_code, crr.nominal, crr.loc_name, crr.rate_value
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
			slog.Warn(err.Error())
		}
		items = append(items, ri)
	}
	return items, nil
}

func (r *repository) LatestRate(ctx context.Context) (*CurrRate, error) {
	query :=
		`SELECT cr.id, cr.for_date, cr.fetch_time, cr.name_val 
		 FROM currency_rate cr 
		 WHERE cr.for_date = (SELECT max(cr1.for_date) FROM currency_rate cr1)`
	row := r.db.QueryRow(ctx, query)

	var cr CurrRate
	if err := row.Scan(&cr.Id, &cr.Date, &cr.FetchDate, &cr.Name); err != nil {
		slog.Warn(err.Error())
		return nil, err
	}

	records, err := r.FindRateRecordsByRateId(ctx, int(cr.Id))
	if err != nil {
		return nil, err
	}
	cr.RateItems = records
	return &cr, nil
}

func (r *repository) IdCurrencyByCharCode(ctx context.Context, code string) (string, error) {
	var curId string
	query := `SELECT id FROM currency_item WHERE iso_char_code=$1`
	if err := r.db.QueryRow(ctx, query, code).Scan(&curId); err != nil {
		return "", err
	}
	return curId, nil
}
