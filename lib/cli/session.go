package cli

import (
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var NoColorsVar = "NOCOLORS"

var ansiColorCodes = regexp.MustCompile(`\x1b\[[^m]+m`)

type Session interface {
	io.Reader
	io.Writer
	Exit(code int) error
	Environ() []string
	Stderr() io.Writer
	Colors() bool
	User() string
}

func ContextSession(cmd *cobra.Command) Session {
	return Context(cmd).Value("session").(Session)
}

type localSession struct{}

func LocalSession() Session {
	return &localSession{}
}

func (ls *localSession) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (ls *localSession) Write(p []byte) (n int, err error) {
	if ls.Colors() {
		return os.Stdout.Write(p)
	}
	return os.Stdout.Write(ansiColorCodes.ReplaceAll(p, []byte{}))
}

func (ls *localSession) Exit(code int) error {
	os.Exit(code)
	return nil
}

func (ls *localSession) Environ() []string {
	return os.Environ()
}

func (ls *localSession) Stderr() io.Writer {
	return os.Stderr
}

func (ls *localSession) Colors() bool {
	return getEnv(ls.Environ(), NoColorsVar) == ""
}

func (ls *localSession) User() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Name
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
