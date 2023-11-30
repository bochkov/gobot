package cbr

import (
	"context"
	"time"

	"github.com/bochkov/gobot/internal/push"
)

type Currency struct {
	//XMLName xml.Name   `xml:"Valuta"`
	//Text    string     `xml:",chardata"`
	Name  string     `xml:"name,attr"`
	Items []CurrItem `xml:"Item"`
}

type CurrItem struct {
	//Text        string `xml:",chardata"`
	Id          string `xml:"ID,attr"`
	Name        string `xml:"Name"`
	EngName     string `xml:"EngName"`
	Nominal     int    `xml:"Nominal"`
	ParentCode  string `xml:"ParentCode"`
	IsoNumCode  int    `xml:"ISO_Num_Code"`
	IsoCharCode string `xml:"ISO_Char_Code"`
}

type CurrRate struct {
	//XMLName   xml.Name   `xml:"ValCurs"`
	//Text      string     `xml:",chardata"`
	Id        int64
	Date      shortDate `xml:"Date,attr"`
	FetchDate time.Time
	Name      string     `xml:"name,attr"`
	RateItems []RateItem `xml:"Valute"`
}

type RateItem struct {
	//Text     string `xml:",chardata"`
	Id       int64
	CurID    string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    float  `xml:"Value"`
}

type CurrRange struct {
	//XMLName    xml.Name `xml:"ValCurs"`
	//Text       string   `xml:",chardata"`
	ID         string            `xml:"ID,attr"`
	DateRange1 shortDate         `xml:"DateRange1,attr"`
	DateRange2 shortDate         `xml:"DateRange2,attr"`
	Name       string            `xml:"name,attr"`
	Records    []CurrRangeRecord `xml:"Record"`
}

type CurrRangeRecord struct {
	//Text    string `xml:",chardata"`
	ID      string    `xml:"Id,attr"`
	Date    shortDate `xml:"Date,attr"`
	Nominal int       `xml:"Nominal"`
	Value   float     `xml:"Value"`
}

type Repository interface {
	IdCurrencyByCharCode(ctx context.Context, code string) (string, error)
	LatestRate(ctx context.Context) (*CurrRate, error)
}

type TaskRepository interface {
	TruncCurrencyItems(ctx context.Context) error
	SaveCurrRate(ctx context.Context, cr CurrRate)
	SaveCurrency(ctx context.Context, c Currency)
}

type Service interface {
	push.Push
	LatestRate(ctx context.Context) (*CurrRate, error)
	LatestRange(ctx context.Context, currencies []string) []CalcRange
	RangeOf(ctx context.Context, code string, from time.Time, to time.Time) (*CurrRange, error)
	Description() string
}

type Handler struct {
	Service
}
