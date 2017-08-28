package inspector

import (
	"encoding/json"
	"net/http"

	"github.com/gliderlabs/cmd/app/console"
	"github.com/gliderlabs/cmd/lib/web"
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("inspector", &Component{})
}

type Component struct{}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return r.URL.Path == "/_inspect"
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := web.Sessions.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := make(map[string]interface{})
	for k, v := range session.Values {
		m[k.(string)] = v
	}
	data := map[string]interface{}{
		"Session": m,
		"User":    console.SessionUser(r),
	}
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.RenderTemplate(w, r, "inspect", map[string]interface{}{
		"JSON": string(b),
	})
}
