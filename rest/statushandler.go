package api

import (
	"github.com/codahale/metrics"
	_ "github.com/codahale/metrics/runtime"
	"github.com/spf13/viper"
	"net/http"
)

func (h *RequestHandler) HandleStatusRequest(w http.ResponseWriter, r *http.Request) error {
	counters, guages := metrics.Snapshot()
	return writeJsonResponse(w, Status{counters, guages, viper.AllSettings()}, nil)
}

type Status struct {
	Counters map[string]uint64
	Guages   map[string]int64
	Config   map[string]interface{}
}
