package builtin

import (
	"io"
	"regexp"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/gliderlabs/cmd/lib/cli"
)

var ansiColorCodes = regexp.MustCompile(`\x1b\[[^m]+m`)

func Execute(sess ssh.Session) error {
	args := sess.Command()
	if len(args) > 0 {
		args[0] = strings.TrimLeft(args[0], ":")
	}
	cliSession := &session{sess}
	root := rootCmd(cliSession)
	root.SetOutput(sess)
	root.SetHelpTemplate(helpTemplate)
	root.SetUsageTemplate(usageTemplate)
	root.SetUsageFunc(rootUsageFunc)
	root.DisableFlagParsing = true
	for _, cmd := range Commands() {
		cli.AddCommand(root, cmd, cliSession)
	}
	return cli.Execute(root, cliSession, args)
}

type session struct {
	ssh.Session
}

func (s *session) Write(p []byte) (n int, err error) {
	if s.Colors() {
		return s.Session.Write(p)
	}
	return s.Session.Write(ansiColorCodes.ReplaceAll(p, []byte{}))
}

func (s *session) Stderr() io.Writer {
	return s.Session.Stderr()
}

func (s *session) Colors() bool {
	return getEnv(s.Environ(), cli.NoColorsVar) == ""
}

func getEnv(environ []string, key string) string {
	for _, envVar := range environ {
		kvp := strings.SplitN(envVar, "=", 2)
		if kvp[0] == key && len(kvp) > 1 {
			return kvp[1]
		}
	}
	return ""
}
