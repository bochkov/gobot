package cbr

import (
	"context"
	"log"
	"sort"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/bochkov/gobot/internal/db"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
)

const (
	CurrencyUrl = "https://www.cbr.ru/scripts/XML_valFull.asp"
	DailyUrl    = "https://cbr.ru/scripts/XML_daily.asp"
	DynamicUrl  = "https://cbr.ru/scripts/XML_dynamic.asp"
)

const (
	insCurrencyItemQuery     = `insert into currency_item (id, name, eng_name, nominal, parent_code, iso_num_code, iso_char_code) values ($1, $2, $3, $4, $5, $6, $7)`
	insCurrencyRateQuery     = `insert into currency_rate (date, name) VALUES ($1, $2) returning id`
	insCurrencyRateItemQuery = `insert into currency_rate_record (curr_id, rate_id, num_code, char_code, nominal, name, rate_value) values ($1, $2, $3, $4, $5, $6, $7)`
)

type Cbr struct {
}

func NewCbr() *Cbr {
	return new(Cbr)
}

func (c *Cbr) fetchAndSaveCurrencies() {
	log.Print("fetch currencies")
	ctx := context.Background()
	var data string
	if err := requests.URL(CurrencyUrl).
		UserAgent("curl/8.0.1").
		ToString(&data).
		Fetch(ctx); err != nil {
		log.Print(err)
		return
	}
	currency := new(Currency)
	if err := util.FromXml(data, &currency); err != nil {
		log.Print(err)
		return
	}
	if _, err := db.GetPool().Exec(ctx, `truncate table currency_item`); err != nil {
		log.Print(err)
		return
	}
	for _, it := range currency.Items {
		if _, err := db.GetPool().Exec(ctx, insCurrencyItemQuery,
			it.Id, it.Name, it.EngName, it.Nominal, it.ParentCode, it.IsoNumCode, it.IsoCharCode); err != nil {
			log.Print(err)
		}
	}
}

func (c *Cbr) fetchAndSaveCurrRates() {
	log.Print("fetch currency rates")
	ctx := context.Background()
	var data string
	if err := requests.URL(DailyUrl).
		UserAgent("curl/8.0.1").
		ToString(&data).
		Fetch(ctx); err != nil {
		log.Print(err)
		return
	}

	currRate := new(CurrRate)
	if err := util.FromXml(data, &currRate); err != nil {
		log.Print(err)
		return
	}
	var id int
	if err := db.GetPool().QueryRow(ctx, insCurrencyRateQuery, currRate.Date, currRate.Name).Scan(&id); err != nil {
		log.Print(err)
		return
	}
	for _, it := range currRate.RateItems {
		if _, err := db.GetPool().Exec(ctx, insCurrencyRateItemQuery,
			it.CurID, id, it.NumCode, it.CharCode, it.Nominal, it.Name, it.Value); err != nil {
			log.Print(err)
		}
	}
}

func (c *Cbr) needUpdate() bool {
	msk, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now()
	year, month, day := now.Date()
	curLoc := now.Location()
	updTime := time.Date(year, month, day, 12, 0, 0, 0, msk)
	midnight := time.Date(year, month, day, 0, 0, 0, 0, curLoc)
	curRate, err := c.findLatestRate()
	if err != nil {
		return true
	}
	return curRate.FetchDate.Before(midnight) ||
		(curRate.Date.Time.Before(time.Now()) && curRate.FetchDate.After(updTime))
}

func (c *Cbr) updateRatesIfNeeded() {
	if c.needUpdate() {
		c.fetchAndSaveCurrRates()
	}
}

func (c *Cbr) findLatestRate() (*CurrRate, error) {
	ctx := context.Background()
	query := `select r.id, r.date, r.fetch_time, r.name from currency_rate r where r.date = (select max(r1.date) from currency_rate r1)`
	var cr CurrRate
	if err := db.GetPool().QueryRow(ctx, query).Scan(&cr.Id, &cr.Date, &cr.FetchDate, &cr.Name); err != nil {
		log.Print(err)
		return nil, err
	}
	query2 := `select r.id, r.curr_id, r.num_code, r.char_code, r.nominal, r.name, r.rate_value from currency_rate_record r where r.rate_id = $1`
	items, err := db.GetPool().Query(ctx, query2, cr.Id)
	defer items.Close()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var cri []RateItem
	for items.Next() {
		var ri RateItem
		if err = items.Scan(&ri.Id, &ri.CurID, &ri.NumCode, &ri.CharCode, &ri.Nominal, &ri.Name, &ri.Value); err != nil {
			log.Print(err)
			return nil, err
		}
		cri = append(cri, ri)
	}
	cr.RateItems = cri
	return &cr, nil
}

func (c *Cbr) rangeOf(code string, from time.Time, to time.Time) (*CurrRange, error) {
	ctx := context.Background()
	var count int
	if err := db.GetPool().QueryRow(ctx, `select count(*) from currency_item`).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		c.fetchAndSaveCurrencies()
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
		log.Print(err)
		return nil, err
	}

	curRange := new(CurrRange)
	if err := util.FromXml(data, &curRange); err != nil {
		return nil, err
	}
	return curRange, nil
}

func (c *Cbr) LatestRange(currencies []string) []CalcRange {
	ranges := make([]CalcRange, 0)
	t := time.Now()
	for _, cur := range currencies {
		r, err := c.rangeOf(cur, t.AddDate(0, 0, -14), t.AddDate(0, 0, 1))
		if err != nil {
			log.Print(err)
		} else {
			sort.Sort(CurrRangeRecordByDateReverse(r.Records))
			rng := NewCalcRange(cur, r.Records[0], r.Records[1])
			ranges = append(ranges, *rng)
		}
	}
	return ranges
}

func (c *Cbr) Text() string {
	var currencies = c.LatestRange([]string{"USD", "EUR"})
	return strings.Join(func() []string {
		var ranges = make([]string, 0)
		for _, cur := range currencies {
			ranges = append(ranges, cur.String())
		}
		return ranges
	}(), "\n")
}
