package console

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"strings"
	"time"
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
	buf := bytes.NewBufferString(code)
	_, err := binary.ReadUvarint(buf)
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
