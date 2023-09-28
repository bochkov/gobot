package anekdot

import (
	"github.com/bochkov/gobot/util"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Anek struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type AnekService interface {
	GetRandom() (*Anek, error)
	GetById(id int) (*Anek, error)
}

func AnekHandler(w http.ResponseWriter, req *http.Request) {
	var baneks AnekService = NewBaneks()
	if req.URL.Query().Get("id") != "" {
		id, _ := strconv.Atoi(req.URL.Query().Get("id"))
		var anek, err = baneks.GetById(id)
		util.Resp(w, anek, err)
	} else if vars := mux.Vars(req)["id"]; vars != "" {
		id, _ := strconv.Atoi(vars)
		var anek, err = baneks.GetById(id)
		util.Resp(w, anek, err)
	} else {
		var anek, err = baneks.GetRandom()
		util.Resp(w, anek, err)
	}
}
