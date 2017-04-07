package sentry

import (
	"encoding/json"
	"net/http"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/events"
	"github.com/gliderlabs/comlab/pkg/log"
)

func (c *Component) MatchHTTP(r *http.Request) bool {
	return r.URL.Path == com.GetString("endpoint")
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event IssueEvent
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&event); err != nil {
		log.Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events.Emit(event)
}
