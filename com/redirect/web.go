package redirect

import "net/http"

func (c *Component) MatchHTTP(r *http.Request) bool {
	return r.URL.Path == "/"
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://gliderlabs.com/devlog/2016/announcing-cmd-io/", http.StatusFound)
}
