package cbr

import (
	"github.com/bochkov/gobot/util"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func LatestRate(w http.ResponseWriter, _ *http.Request) {
	var cbr = NewCbr()
	cbr.updateRatesIfNeeded()
	curRate, err := cbr.findLatestRate()
	util.Resp(w, curRate, err)
}

func LatestRates(w http.ResponseWriter, req *http.Request) {
	var cbr = NewCbr()
	cbr.updateRatesIfNeeded()
	currency := req.URL.Query().Get("currency")
	if currency == "" {
		currency = "USD,EUR"
	}
	ranges := cbr.LatestRange(strings.Split(strings.ToUpper(currency), ","))
	util.Resp(w, ranges, nil)
}

func PeriodRates(w http.ResponseWriter, req *http.Request) {
	var cbr = NewCbr()
	cbr.updateRatesIfNeeded()
	vars := mux.Vars(req)
	t := time.Now()
	currency := strings.ToUpper(vars["currency"])
	switch vars["period"] {
	case "month":
		r, err := cbr.rangeOf(currency, t.AddDate(0, -1, 0), t.AddDate(0, 0, 1))
		util.Resp(w, r, err)
	case "year":
		r, err := cbr.rangeOf(currency, t.AddDate(-1, 0, 0), t.AddDate(0, 0, 1))
		util.Resp(w, r, err)
	}
}
