package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
	uuid "github.com/satori/go.uuid"

	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
	"github.com/progrium/cmd/pkg/access"
)

const rootUsageTmpl = `Usage:{{if .Runnable}}
  ssh <user>@{{.UseLine}}{{ if .HasAvailableSubCommands}} [command]{{end}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "[command] --help" for help about a meta command.{{end}}

`

const metaUsageTmpl = `Usage:{{if .Runnable}}{{if not .HasParent }}
  ssh <user>@cmd.io {{.UseLine}}{{ if .HasAvailableSubCommands}}:[command]{{end}}{{else}}
  {{.UseLine}}{{ if .HasAvailableSubCommands}} [command]{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "[command] --help" for help about a meta command.{{end}}

`

func HandleSSH(s ssh.Session) {
	var (
		start = time.Now()
		msg   = ""
		cmd   = &core.Command{}
		user  = s.User()
	)
	defer func() {
		log.Info(s, cmd, time.Since(start), msg)
	}()

	// check for channel access when user is not a token
	if tok := uuid.FromStringOrNil(user); tok == uuid.Nil && !access.Check(user) {
		msg = "channel access denied"
		fmt.Fprintln(s, com.GetString("access_denied_msg"))
		return
	}

	args := s.Command()
	if len(args) == 0 {
		args = []string{":help"}
	}

	if strings.HasPrefix(args[0], ":") {
		runRootMeta(s, args)
		return
	}

	stream := &core.Stream{
		Stdin:  s,
		Stdout: s,
		Stderr: s.Stderr(),
	}
	// handle git-receive-pack by finding the first cmd which has io.cmd.git-receive == arg[1]
	if strings.HasPrefix(args[0], "git-receive-pack") && len(args) > 1 {
		cmds := store.Selected().List(user)
		args[1] = strings.TrimPrefix(args[1], "/")
		for _, c := range cmds {
			path, ok := c.Environment["io.cmd.git-receive"]
			if ok && strings.HasPrefix(args[1], path) {
				cmd = c
				c.Run(stream, s.User(), args)
				return
			}
		}
	}

	if strings.Contains(args[0], ":") {
		parts := strings.Split(args[0], ":")
		if parts[1] == "" {
			parts[1] = "help"
		}
		args[0] = ":" + parts[1]
		if parts[0] == "" {
			msg = "command missing"
			fmt.Fprintln(s.Stderr(), "Command name required for meta command: "+args[0])
			s.Exit(1)
			return
		}
		runCmdMeta(s, parts[0], args)
		return
	}

	if strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		cmd = store.Selected().Get(parts[0], parts[1])
		if cmd == nil {
			msg = "command not found"
			fmt.Fprintln(s.Stderr(), "Command not found: "+args[0])
			s.Exit(1)
			return
		}
		if !cmd.HasAccess(user) {
			msg = "cmd access denied"
			fmt.Fprintln(s.Stderr(), "Not allowed")
			s.Exit(1)
			return
		}

		cmd.Run(stream, s.User(), args[1:])
		return
	}

	cmd = store.Selected().Get(user, args[0])
	if cmd == nil {
		if cmd = LazyLoad(user, args[0]); cmd == nil {
			msg = "command not found"
			fmt.Fprintln(s.Stderr(), "Command not found: "+args[0])
			s.Exit(1)
			return
		}
		if err := store.Selected().Put(user, args[0], cmd); err != nil {
			fmt.Fprintln(s.Stderr(), err.Error())
			s.Exit(255)
			return
		}
	}
	cmd.Run(stream, s.User(), args[1:])
	if cmd.Changed {
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(s.Stderr(), err.Error())
			s.Exit(255)
			return
		}
	}
}

func runRootMeta(s ssh.Session, args []string) {
	root := &MetaCommand{Use: "cmd.io", Session: s}
	root.Add(RootCommands()...)
	root.Cmd.SetArgs(args)
	root.Cmd.SetOutput(s)
	root.Cmd.SetUsageTemplate(rootUsageTmpl)
	if err := root.Cmd.Execute(); err != nil {
		s.Exit(255)
	}
}

func runCmdMeta(s ssh.Session, cmdName string, args []string) {
	var cmd *core.Command
	if strings.Contains(cmdName, "/") {
		parts := strings.SplitN(cmdName, "/", 2)
		cmd = store.Selected().Get(parts[0], parts[1])
	} else {
		cmd = store.Selected().Get(s.User(), cmdName)
	}
	if cmd == nil {
		fmt.Fprintln(s.Stderr(), "Command not found: "+cmdName)
		s.Exit(1)
		return
	}
	if !cmd.IsAdmin(s.User()) {
		fmt.Fprintln(s.Stderr(), "Not allowed")
		s.Exit(1)
		return
	}
	meta := &MetaCommand{Use: cmdName, Session: s, ForCmd: cmd}
	meta.Add(MetaCommands()...)
	meta.Cmd.SetArgs(args)
	meta.Cmd.SetOutput(s)
	meta.Cmd.SetUsageTemplate(metaUsageTmpl)
	if err := meta.Cmd.Execute(); err != nil {
		log.Debug(s, err)
		s.Exit(255)
	}
}

func HandleAuth(user string, key ssh.PublicKey) bool {
	if tok := uuid.FromStringOrNil(user); tok != uuid.Nil {
		token, _ := store.Selected().GetToken(tok.String())
		if token != nil && token.Key == user {
			return true
		}
		log.Info("no match found for token: " + user)
	}

	resp, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", user))
	if err != nil {
		log.Info(user, err)
		return false
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		k, _, _, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
		if err != nil {
			log.Info(user, err)
			continue
		}

		if ssh.KeysEqual(key, k) {
			return true
		}
	}
	log.Info("auth: no matching keys for: "+user, key)
	return false
}

func LazyLoad(user, name string) *core.Command {
	cmd := &core.Command{
		Name:   name,
		User:   user,
		Source: fmt.Sprintf("%s/cmd-%s", user, name),
	}
	if err := cmd.Pull(); err != nil {
		log.Debug(cmd, err)
		return nil
	}
	return cmd
}
