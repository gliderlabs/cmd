package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gliderlabs/pkg/com"
	log "github.com/gliderlabs/pkg/log"
)

type Allowed struct {
	last time.Time
	ttl  time.Duration

	users map[string]bool
	mu    sync.Mutex // protects users map
}

func (a *Allowed) Check(name string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.users == nil {
		a.users = make(map[string]bool)
	}
	if time.Since(a.last) < a.ttl {
		return a.users[name]
	}
	a.last = time.Now()
	url := fmt.Sprintf("https://api.github.com/teams/%s/members?access_token=%s",
		com.GetString("gh_team_id"), com.GetString("gh_token"))

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "cmd.io")
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Info(err)
		return a.users[name]
	}

	members := []struct {
		Login string
	}{}
	json.NewDecoder(res.Body).Decode(&members)
	defer res.Body.Close()
	for _, m := range members {
		a.users[m.Login] = true
	}

	return a.users[name]
}
