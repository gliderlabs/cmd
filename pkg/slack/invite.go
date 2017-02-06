package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// InviteToTeam a user by email using a token with admin privileges
func InviteToTeam(team, token, email string) error {
	resp, err := http.PostForm(
		fmt.Sprintf("https://%s.slack.com/api/users.admin.invite", team),
		url.Values{"email": {email}, "token": {token}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var msg struct {
		Ok    bool
		Error string
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}
	if !msg.Ok {
		return errors.New(msg.Error)
	}
	return nil
}
