package cbr

import (
	"context"
	"sort"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
	"log/slog"
)

const (
	CurrencyUrl = "https://www.cbr.ru/scripts/XML_valFull.asp"
	DailyUrl    = "https://cbr.ru/scripts/XML_daily.asp"
	DynamicUrl  = "https://cbr.ru/scripts/XML_dynamic.asp"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(3) * time.Second,
	}
}

func (s *service) Description() string {
	return "Официальный курс валют"
}

func (s *service) LatestRate(c context.Context) (*CurrRate, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// TODO s.updateRatesIfNeeded()
	cr, err := s.Repository.LatestRate(ctx)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (s *service) LatestRange(c context.Context, currencies []string) []CalcRange {
	// TODO s.updateRatesIfNeeded()
	ranges := make([]CalcRange, 0)
	t := time.Now()
	for _, cur := range currencies {
		r, err := s.RangeOf(c, cur, t.AddDate(0, 0, -14), t.AddDate(0, 0, 1))
		if err != nil {
			slog.Warn(err.Error())
		} else {
			sort.Sort(CurrRangeRecordByDateReverse(r.Records))
			rng := NewCalcRange(cur, r.Records[0], r.Records[1])
			ranges = append(ranges, *rng)
		}
	}
	return ranges
}

func (s *service) RangeOf(c context.Context, code string, from time.Time, to time.Time) (*CurrRange, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// TODO s.updateRatesIfNeeded()

	var count int
	if err := db.GetPool().QueryRow(ctx, `select count(*) from currency_item`).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		s.fetchAndSaveCurrencies()
	}
	var curId string
	if err := db.GetPool().QueryRow(ctx, `select id from currency_item where iso_char_code=$1`, code).Scan(&curId); err != nil {
		return nil, err
	}
	var data string
	if err := requests.URL(DynamicUrl).
		UserAgent("curl/8.0.1").
		Param("date_req1", from.Format("02/01/2006")).
		Param("date_req2", to.Format("02/01/2006")).
		Param("VAL_NM_RQ", curId).
		ToString(&data).
		Fetch(ctx); err != nil {
		slog.Warn(err.Error())
		return nil, err
	}

	curRange := new(CurrRange)
	if err := util.FromXml(data, &curRange); err != nil {
		return nil, err
	}
	return curRange, nil
}

func (s *service) fetchAndSaveCurrencies() {
	slog.Debug("fetch currencies")
	ctx := context.Background()
	var data string
	if err := requests.URL(CurrencyUrl).
		UserAgent("curl/8.0.1").
		ToString(&data).
		Fetch(ctx); err != nil {
		slog.Warn(err.Error())
		return
	}
	currency := new(Currency)
	if err := util.FromXml(data, &currency); err != nil {
		slog.Warn(err.Error())
		return
	}
	if _, err := db.GetPool().Exec(ctx, `truncate table currency_item`); err != nil {
		slog.Warn(err.Error())
		return
	}
	for _, it := range currency.Items {
		if _, err := db.GetPool().Exec(ctx, insCurrencyItemQuery,
			it.Id, it.Name, it.EngName, it.Nominal, it.ParentCode, it.IsoNumCode, it.IsoCharCode); err != nil {
			slog.Warn(err.Error())
		}
	}
}

func (s *service) needUpdate() bool {
	msk, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now()
	year, month, day := now.Date()
	curLoc := now.Location()
	updTime := time.Date(year, month, day, 12, 0, 0, 0, msk)
	midnight := time.Date(year, month, day, 0, 0, 0, 0, curLoc)
	curRate, err := s.LatestRate(context.Background())
	if err != nil {
		return true
	}
	return curRate.FetchDate.Before(midnight) ||
		(curRate.Date.Time.Before(time.Now()) && curRate.FetchDate.After(updTime))
}

func (s *service) updateRatesIfNeeded() {
	if s.needUpdate() {
		slog.Debug("fetch currency rates")
		ctx := context.Background()
		var data string
		if err := requests.URL(DailyUrl).
			UserAgent("curl/8.0.1").
			ToString(&data).
			Fetch(ctx); err != nil {
			slog.Warn(err.Error())
			return
		}

		var currRate CurrRate
		if err := util.FromXml(data, &currRate); err != nil {
			slog.Warn(err.Error())
			return
		}
		var id int
		if err := db.GetPool().QueryRow(ctx, insCurrencyRateQuery, currRate.Date, currRate.Name).Scan(&id); err != nil {
			slog.Warn(err.Error())
			return
		}
		for _, it := range currRate.RateItems {
			if _, err := db.GetPool().Exec(ctx, insCurrencyRateItemQuery,
				it.CurID, id, it.NumCode, it.CharCode, it.Nominal, it.Name, it.Value); err != nil {
				slog.Warn(err.Error())
			}
		}
	}
}

func (s *service) Text() string {
	var currencies = s.LatestRange(context.Background(), []string{"USD", "EUR"})
	return strings.Join(func() []string {
		var ranges = make([]string, 0)
		for _, cur := range currencies {
			ranges = append(ranges, cur.String())
		}
		return ranges
	}(), "\n")
}
