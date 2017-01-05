package cmd

import (
	"strconv"

	"golang.org/x/oauth2"

	"github.com/gliderlabs/gosper/pkg/com"
	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
)

type Allowed struct {
	gh *github.Client
}

func (a *Allowed) client() *github.Client {
	if a.gh == nil {
		auth := &oauth2.Transport{Source: oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: com.GetString("gh_token")},
		)}
		cache := httpcache.NewMemoryCacheTransport()
		cache.Transport = auth
		a.gh = github.NewClient(cache.Client())
		a.gh.UserAgent = "cmd.io"
	}
	return a.gh
}

func (a *Allowed) Check(name string) bool {
	id, _ := strconv.Atoi(com.GetString("gh_team_id"))
	isMember, res, err := a.client().Organizations.IsTeamMember(id, name)
	if err != nil {
		log.Info(err)
		return isMember
	}
	log.Info("github api rate:", res.Rate.String())
	if res.Header.Get(httpcache.XFromCache) != "" {
		log.Info("member: " + name + " checked from cache")
	}
	return isMember
}
