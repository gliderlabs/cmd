package console

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/cmd/pkg/auth0"
	"github.com/gliderlabs/comlab/pkg/log"
)

const (
	MaxPendingInvites = 5
)

func GenerateInviteCode(user *User) (string, error) {
	id := strings.Split(user.ID, "|")
	i, err := strconv.Atoi(id[1])
	if err != nil {
		return "", err
	}
	buf := make([]byte, 12)
	n := binary.PutUvarint(buf, uint64(time.Now().UnixNano()))
	nn := binary.PutUvarint(buf[n:], uint64(i))
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf[:n+nn]), nil
}

func InviteCodeUser(code string) (*User, error) {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(code)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	_, err = binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	id, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	uid := "github|" + strconv.Itoa(int(id))
	user, err := LookupUser(uid)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func codesHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pendingInvites := user.Account.Invites.Pending
	if r.Method == "POST" {
		if len(user.Account.Invites.Pending) <= MaxPendingInvites {
			code, err := GenerateInviteCode(user)
			if err != nil {
				log.Info(r, err, log.Fields{"uid": user.ID})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pendingInvites = append(pendingInvites, code)
			err = auth0.DefaultClient().PatchUser(user.ID, auth0.User{
				"app_metadata": map[string]interface{}{
					"invites": map[string]interface{}{
						"pending":    pendingInvites,
						"invited_by": user.Account.Invites.InvitedBy,
					},
				},
			})
			if err != nil {
				log.Info(r, err, log.Fields{"uid": user.ID})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	w.Header().Add("content-type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]interface{}{
		"pending": pendingInvites,
		"max":     MaxPendingInvites,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func inviteHandler(w http.ResponseWriter, r *http.Request) {
	var code string
	parts := strings.Split(r.URL.Path, "/")
	if parts[len(parts)-1] == "" {
		code = parts[len(parts)-2]
	} else {
		code = parts[len(parts)-1]
	}
	http.Redirect(w, r, fmt.Sprintf("/register?code=%s", code), http.StatusFound)
}
