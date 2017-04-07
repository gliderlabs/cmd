package slack

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/nlopes/slack"
)

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, com.GetString("endpoint"))
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: mux for other endpoints
	if !strings.HasSuffix(r.URL.Path, "/msg") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.URL.Query().Get("token") != com.GetString("token") {
		http.Error(w, "Bad token", http.StatusForbidden)
		return
	}
	var msg slack.Attachment
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := PostAttachment(msg, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
