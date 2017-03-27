package stripe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
)

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		// mainly for pub_key
		"stripe": func(key string) string {
			return com.GetString(key)
		},
	}
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return r.URL.Path == com.GetString("event_endpoint")
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventID, eventType, err := parseEvent(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listeners := EventListeners(eventType)
	if len(listeners) > 0 {
		event, err := Client().Events.Get(eventID, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, listener := range listeners {
			listener.StripeReceive(*event)
		}
	}
}

func parseEvent(r *http.Request) (string, string, error) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", "", err
	}
	var event map[string]interface{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		return "", "", err
	}
	id, ok := event["id"].(string)
	if !ok {
		return "", "", fmt.Errorf("no id field in event")
	}
	typ, ok := event["type"].(string)
	if !ok {
		return "", "", fmt.Errorf("no type field in event")
	}
	return id, typ, nil
}
