package auto

import (
	"github.com/bochkov/gobot/internal/util"
	"net/http"
)

func CodesHandler(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	region, err := FindRegionByCode(code)
	util.Resp(w, region, err)
}

func RegionsHandler(w http.ResponseWriter, req *http.Request) {
	region := req.URL.Query().Get("region")
	if region == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	codes, err := FindRegionByName(region)
	util.Resp(w, codes, err)
}
