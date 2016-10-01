package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/pkg/com"
	"github.com/gliderlabs/pkg/com/viper"
	"github.com/gliderlabs/pkg/log"
	"github.com/gliderlabs/pkg/ssh"
	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/thejerf/suture"
)

var Store CommandStore

var allowedUsers = &Allowed{}

func LocalMode() bool {
	return os.Getenv("LOCAL") != "false"
}

func DebugMode() bool {
	return os.Getenv("DEBUG") != ""
}

func Run() {
	log.RegisterObserver(new(logging))
	log.SetFieldProcessor(fieldProcessor)

	cfg := viper.NewConfig()
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	com.SetConfig(cfg)

	libhoney.Init(libhoney.Config{
		WriteKey:   com.GetString("honeycomb_key"),
		Dataset:    com.GetString("honeycomb_dataset"),
		SampleRate: 1,
	})
	// when all done, call close
	defer libhoney.Close()
	hostname, _ := os.Hostname()
	libhoney.AddField("servername", hostname)
	libhoney.AddField("release", os.Getenv("RELEASE"))

	log.RegisterObserver(new(honeylog))

	log.RegisterObserver(newRavenLogger(com.GetString("sentry_dsn")))

	Store = GetDynamodbStore()
	app := suture.NewSimple("cmd.io")
	for _, service := range com.Enabled(new(suture.Service), nil) {
		app.Add(service.(suture.Service))
	}
	app.Serve()
}

func HandleSSH(s ssh.Session) {
	var start, cmd = time.Now(), &Command{}
	defer func() {
		log.Info(s, cmd, time.Since(start))
	}()
	user := s.User()
	if !allowedUsers.Check(user) {
		fmt.Fprintln(s, com.GetString("access_denied_msg"))
		return
	}
	args := s.Command()
	if len(args) == 0 {
		args = []string{":"}
	}

	if strings.Contains(args[0], ":") {
		parts := strings.Split(args[0], ":")
		if parts[1] == "" {
			parts[1] = "help"
		}
		args[0] = ":" + parts[1]
		if parts[0] == "" {
			runRootMeta(s, args)
			return
		}
		runCmdMeta(s, parts[0], args)
		return
	}

	if strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		cmd = Store.Get(parts[0], parts[1])
		if cmd == nil {
			fmt.Fprintln(s.Stderr(), "Command not found: "+args[0])
			s.Exit(1)
			return
		}
		if !cmd.HasAccess(user) {
			fmt.Fprintln(s.Stderr(), "Not allowed")
			s.Exit(1)
			return
		}
		cmd.Run(s, args[1:])
		return
	}

	cmd = Store.Get(user, args[0])
	if cmd == nil {
		if cmd = LazyLoad(user, args[0]); cmd == nil {
			fmt.Fprintln(s.Stderr(), "Command not found: "+args[0])
			s.Exit(1)
			return
		}
		if err := Store.Put(user, args[0], cmd); err != nil {
			fmt.Fprintln(s.Stderr(), err.Error())
			s.Exit(255)
			return
		}
	}
	cmd.Run(s, args[1:])
}

func runRootMeta(s ssh.Session, args []string) {
	root := &MetaCommand{Use: "cmd.io", Session: s}
	root.Add(rootHelp, rootInstall, rootUninstall, rootList)
	root.Cmd.SetArgs(args)
	root.Cmd.SetOutput(s)
	root.Cmd.SetUsageTemplate(rootUsageTmpl)
	if err := root.Cmd.Execute(); err != nil {
		//fmt.Fprintln(s.Stderr(), err.Error())
		s.Exit(255)
	}
}

func runCmdMeta(s ssh.Session, cmdName string, args []string) {
	var cmd *Command
	if strings.Contains(cmdName, "/") {
		parts := strings.SplitN(cmdName, "/", 2)
		cmd = Store.Get(parts[0], parts[1])
	} else {
		cmd = Store.Get(s.User(), cmdName)
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
	configCmd := meta.Add(metaConfig)
	configCmd.Add(metaConfigSet, metaConfigUnset)
	accessCmd := meta.Add(metaAccess)
	if cmd.IsPublic() {
		accessCmd.Add(metaAccessPrivate)
	} else {
		accessCmd.Add(metaAccessPublic, metaAccessAdd, metaAccessRemove)
	}
	adminsCmd := meta.Add(metaAdmins)
	adminsCmd.Add(metaAdminsAdd, metaAdminsRemove)
	meta.Add(metaHelp)
	meta.Cmd.SetArgs(args)
	meta.Cmd.SetOutput(s)
	meta.Cmd.SetUsageTemplate(metaUsageTmpl)
	if err := meta.Cmd.Execute(); err != nil {

		//fmt.Fprintln(s.Stderr(), err.Error())
		log.Debug(s, err)
		s.Exit(255)
	}
}

func HandleAuth(user string, key ssh.PublicKey) bool {
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
			return false
		}

		if ssh.KeysEqual(key, k) {
			return true
		}
	}
	return false
}

func LazyLoad(user, name string) *Command {
	cmd := &Command{
		Name:   name,
		User:   user,
		Config: make(map[string]string),
		Source: fmt.Sprintf("%s/cmd-%s", user, name),
	}
	if err := cmd.Pull(); err != nil {
		log.Info(cmd, err)
		return nil
	}
	return cmd
}
