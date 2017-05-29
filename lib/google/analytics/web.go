package analytics

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
)

const trackingScript = `<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', '%s', 'auto');
  ga('send', 'pageview');
</script>`

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"googleAnalytics": func() string {
			if com.GetString("tracking_id") == "" {
				return ""
			}
			return fmt.Sprintf(trackingScript, com.GetString("tracking_id"))
		},
	}
}
