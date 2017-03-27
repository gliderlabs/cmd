package access

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
)

func init() {
	com.Register("access", &Component{},
		com.Option("gh_team_id", "2144066", "GitHub team ID to allow access to"),
		com.Option("gh_token", "", "GitHub access token"))
}

type Component struct {
	teamID int
	client *github.Client
}

// TODO: this component is unnecessary, but is around from first pass
// at consolidating into single package from two
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
