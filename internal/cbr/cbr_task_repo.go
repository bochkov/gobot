package cbr

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) TaskRepository {
	return &taskRepo{db: db}
}

func (t *taskRepo) TruncCurrencyItems(ctx context.Context) error {
	_, err := t.db.Exec(ctx, `TRUNCATE TABLE currency_item`)
	return err
}

func (t *taskRepo) SaveCurrRate(ctx context.Context, cr CurrRate) {
	query :=
		`INSERT INTO currency_rate (for_date, name_val) 
		 VALUES ($1, $2) 
		 ON CONFLICT DO NOTHING 
		 RETURNING id`
	query2 :=
		`INSERT INTO currency_rate_record (curr_id, rate_id, num_code, char_code, nominal, loc_name, rate_value) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	t.runInTx(ctx, "saveCurrRate", func(tx pgx.Tx) error {
		var id int
		if err := tx.QueryRow(ctx, query, cr.Date, cr.Name).Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				slog.Error(err.Error())
			} else {
				return err
			}
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}

		for _, it := range cr.RateItems {
			if _, err := t.db.Exec(ctx, query2,
				it.CurID, id, it.NumCode, it.CharCode, it.Nominal, it.Name, it.Value); err != nil {
				slog.Warn(err.Error())
			}
		}
		return nil
	})
}

func (t *taskRepo) SaveCurrency(ctx context.Context, c Currency) {
	query :=
		`INSERT INTO currency_item (id, loc_name, eng_name, nominal, parent_code, iso_num_code, iso_char_code) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	t.runInTx(ctx, "saveCurrency", func(tx pgx.Tx) error {
		for _, it := range c.Items {
			if _, err := tx.Exec(ctx, query,
				it.Id, it.Name, it.EngName, it.Nominal, it.ParentCode, it.IsoNumCode, it.IsoCharCode); err != nil {
				slog.Warn(err.Error())
			}
		}
		return nil
	})
}

func (t *taskRepo) runInTx(ctx context.Context, desc string, fun func(tx pgx.Tx) error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		slog.Warn("cannot start transaction", "op", desc, "err", err)
		return
	}

	if err := fun(tx); err != nil {
		if er := tx.Rollback(ctx); er != nil {
			slog.Error("cannot rollback transaction", "op", desc, "err", er)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("cannot commit transaction", "op", desc, "err", err)
		}
	}
}
