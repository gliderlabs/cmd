package github

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
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var err error
	switch r.Header.Get("x-github-event") {
	case "issues":
		var event IssuesEvent
		err = dec.Decode(&event)
		if err == nil {
			events.Emit(event)
		}
	case "pull_request":
		var event PullRequestEvent
		err = dec.Decode(&event)
		if err == nil {
			events.Emit(event)
		}
	case "status":
		var event StatusEvent
		err = dec.Decode(&event)
		if err == nil {
			events.Emit(event)
		}
	case "ping":
		var event PingEvent
		err = dec.Decode(&event)
		if err == nil {
			events.Emit(event)
		}
	}
	if err != nil {
		log.Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
