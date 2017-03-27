package access

import (
	"strconv"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"golang.org/x/oauth2"
)

func Check(name string) bool {
	auth := &oauth2.Transport{Source: oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: com.GetString("gh_token")},
	)}
	cache := httpcache.NewMemoryCacheTransport()
	cache.Transport = auth
	gh := github.NewClient(cache.Client())
	gh.UserAgent = "cmd.io"
	id, _ := strconv.Atoi(com.GetString("gh_team_id"))
	checker := &Component{
		teamID: id,
		client: gh,
	}
	return checker.Check(name)
}
