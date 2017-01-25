package stripe

import (
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
