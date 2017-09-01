package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"

	"github.com/gliderlabs/cmd/app/builtin"
	"github.com/gliderlabs/cmd/app/console"
	"github.com/gliderlabs/cmd/app/core"
	"github.com/gliderlabs/cmd/app/store"
)

// Default expiry of 30 sec and expiry purge every 5 min.
// Would be nice to find a good cache with size limit as well.
var authCache = cache.New(30*time.Second, 5*time.Minute)

// TODO: make this more integrated with console?
type cachedUser struct {
	user console.User
	keys []ssh.PublicKey
}

func (c *Component) HandleSSH(s ssh.Session) {
	var (
		start    = time.Now()
		msg      = ""
		cmd      = &core.Command{}
		userName = s.User()
		cmdName  = ""
	)
	defer func() {
		log.Info(s, cmd, time.Since(start), msg, log.Fields{"docker": cmd.Docker().Host})
	}()

	var cont bool
	for _, preprocessor := range Preprocessors() {
		cont, msg = preprocessor.PreprocessSession(s)
		if !cont {
			authCache.Delete(s.User())
			return
		}
	}

	// TODO: make builtins also a preprocessor!
	args := s.Command()
	if len(args) == 0 || strings.HasPrefix(args[0], ":") {
		if err := builtin.Execute(s); err != nil {
			s.Exit(255)
		}
		return
	}
	cmdName = args[0]

	// TODO: make into preprocessor
	// handle git-receive-pack by finding the first cmd which has io.cmd.git-receive == arg[1]
	if strings.HasPrefix(cmdName, "git-receive-pack") && len(args) > 1 {
		cmds := store.Selected().List(userName)
		args[1] = strings.TrimPrefix(args[1], "/")
		for _, c := range cmds {
			path, ok := c.Environment["io.cmd.git-receive"]
			if ok && strings.HasPrefix(args[1], path) {
				cmd = c
				c.Run(s, args)
				return
			}
		}
	}

	if strings.Contains(cmdName, "/") {
		parts := strings.SplitN(cmdName, "/", 2)
		userName = parts[0]
		cmdName = parts[1]
	}

	cmd = store.Selected().Get(userName, cmdName)
	if cmd == nil {
		msg = "command not found"
		fmt.Fprintln(s.Stderr(), "Command not found:", args[0])
		s.Exit(1)
		return
	}
	if !cmd.HasAccess(userName) {
		msg = "cmd access denied"
		fmt.Fprintln(s.Stderr(), "Not allowed")
		s.Exit(1)
		return
	}
	cmd.Run(s, args[1:])
}

func (c *Component) HandleAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	user := ctx.User()
	if tok := uuid.FromStringOrNil(user); tok != uuid.Nil {
		token, _ := store.Selected().GetToken(tok.String())
		if token != nil && token.Key == user {
			return true
		}
		log.Info("no match found for token: " + user)
	}

	var u cachedUser
	cu, ok := authCache.Get(user)
	if ok {
		u = cu.(cachedUser)
	} else {
		resp, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", user))
		if err != nil {
			log.Info(user, err)
			return false
		}
		if resp.StatusCode == http.StatusNotFound {
			log.Info(fmt.Sprintf("github user '%s' not found", user), key)
			return false
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)
		var keys []ssh.PublicKey
		for scanner.Scan() {
			k, _, _, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Info(user, err)
				continue
			}
			keys = append(keys, k)
		}
		usr, err := console.LookupNickname(user)
		if err != nil {
			log.Info(user, err)
		}
		u = cachedUser{
			user: usr,
			keys: keys,
		}
		authCache.Set(user, u, cache.DefaultExpiration)
	}
	ctx.SetValue("user", &(u.user))
	ctx.SetValue("plan", u.user.Account.Plan)

	for _, k := range u.keys {
		if ssh.KeysEqual(key, k) {
			return true
		}
	}
	log.Info(fmt.Sprintf("no matching keys of %d for '%s'", len(u.keys), user), key)
	return false
}
