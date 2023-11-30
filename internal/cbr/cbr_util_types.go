package cbr

import (
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type float float32

func (f *float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}
	raw = strings.ReplaceAll(raw, ",", ".")
	f64, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return err
	}
	*f = float(f64)
	return nil
}

func (f *float) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.4f", *f)
	return []byte(s), nil
}

func (f *float) Abs() float32 {
	f32 := float32(*f)
	return math.Float32frombits(math.Float32bits(f32) &^ (1 << 31))
}

type shortDate struct{ pgtype.Date }

const shortDateForm = "02.01.2006"

func (sd *shortDate) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse(shortDateForm, attr.Value)
	if err != nil {
		return err
	}
	shor := new(shortDate)
	err = shor.Scan(parse)
	if err != nil {
		return err
	}
	*sd = *shor
	return nil
}

type CalcRange struct {
	Code  string          `json:"code"`
	Value CurrRangeRecord `json:"value"`
	Delta float           `json:"delta"`
}

func NewCalcRange(code string, rec0 CurrRangeRecord, rec1 CurrRangeRecord) *CalcRange {
	return &CalcRange{
		Code:  code,
		Value: rec0,
		Delta: rec0.Value - rec1.Value,
	}
}

func (cr *CalcRange) String() string {
	sign := "+"
	if cr.Delta < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%d %s = %.3f ₽ [ %s%.3f ]",
		cr.Value.Nominal, strings.ToUpper(cr.Code), cr.Value.Value, sign, cr.Delta.Abs())
}

type CurrRangeRecordByDateReverse []CurrRangeRecord

func (cmp CurrRangeRecordByDateReverse) Len() int {
	return len(cmp)
}

func (cmp CurrRangeRecordByDateReverse) Less(i, j int) bool {
	return cmp[i].Date.Time.After(cmp[j].Date.Time)
}

func (cmp CurrRangeRecordByDateReverse) Swap(i, j int) {
	cmp[i], cmp[j] = cmp[j], cmp[i]
}
