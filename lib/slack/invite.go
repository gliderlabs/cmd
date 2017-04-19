package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gliderlabs/comlab/pkg/com"
)

// InviteToTeam a user by email
func InviteToTeam(team, email string) error {
	resp, err := http.PostForm(
		fmt.Sprintf("https://%s.slack.com/api/users.admin.invite", team),
		url.Values{"email": {email}, "token": {com.GetString("token")}})
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
