package access

import (
	"strconv"

	"golang.org/x/oauth2"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/progrium/cmd/pkg/access"
)

func Register() {
	com.Register("access", &Component{},
		com.Option("gh_team_id", "2144066", "GitHub team ID to allow access to"),
		com.Option("gh_token", "", "GitHub access token"))
	access.DefaultClientFactory = ClientFactory
}

func ClientFactory() access.Client {
	auth := &oauth2.Transport{Source: oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: com.GetString("gh_token")},
	)}
	cache := httpcache.NewMemoryCacheTransport()
	cache.Transport = auth
	gh := github.NewClient(cache.Client())
	gh.UserAgent = "cmd.io"

	id, _ := strconv.Atoi(com.GetString("gh_team_id"))
	return &Component{
		teamID: id,
		client: gh,
	}
}

type Component struct {
	teamID int
	client *github.Client
}

func (c *Component) Check(name string) bool {
	isMember, res, err := c.client.Organizations.IsTeamMember(c.teamID, name)
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
