package cbr

import (
	"context"
	"log/slog"
	"sort"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
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

	cr, err := s.Repository.LatestRate(ctx)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (s *service) LatestRange(c context.Context, currencies []string) []CalcRange {
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

	curId, err := s.Repository.IdCurrencyByCharCode(ctx, code)
	if err != nil {
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

	var curRange *CurrRange
	if err := util.FromXml(data, &curRange); err != nil {
		return nil, err
	}
	return curRange, nil
}

func (s *service) PushText() string {
	var currencies = s.LatestRange(context.Background(), []string{"USD", "EUR"})
	return strings.Join(func() []string {
		var ranges = make([]string, 0)
		for _, cur := range currencies {
			ranges = append(ranges, cur.String())
		}
		return ranges
	}(), "\n")
}
