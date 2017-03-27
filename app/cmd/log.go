package cmd

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cast"

	"github.com/progrium/cmd/app/core"
)

func fieldProcessor(e log.Event, o interface{}) (log.Event, bool) {
	switch obj := o.(type) {
	case time.Duration:
		return e.Append("dur", cast.ToString(int64(obj/time.Millisecond))), true
	case log.ResponseWriter:
		e = e.Append("bytes", cast.ToString(obj.Size()))
		e = e.Append("status", cast.ToString(obj.Status()))
		return e, true
	case *http.Request:
		e = e.Append("ip", obj.RemoteAddr)
		e = e.Append("method", obj.Method)
		e = e.Append("path", obj.RequestURI)
		return e, true
	case *core.Command:
		if obj == nil {
			return e, false
		}
		e = e.Append("cmd.user", obj.User)
		e = e.Append("cmd.name", obj.Name)
		return e, true
	case ssh.Session:
		e = e.Append("sess.user", obj.User())
		e = e.Append("sess.remoteaddr", obj.RemoteAddr().String())
		e = e.Append("sess.command", strings.Join(obj.Command(), " "))
		return e, true
	case ssh.PublicKey:
		e = e.Append("pubkey.type", obj.Type())
		e = e.Append("pubkey.hash", fmt.Sprintf("%x", md5.Sum(obj.Marshal())))
		return e, true
	}
	return e, false
}
