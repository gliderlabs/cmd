package dev

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/mattaitchison/bashenv"
	"github.com/rjeczalik/notify"
	"github.com/thejerf/suture"
)

var (
	buildErrored = false
)

func Run() {
	cwd, _ := os.Getwd()
	cmdName := filepath.Base(cwd)

	// find entrypoint cmd
	if _, err := os.Stat("cmd/" + cmdName); err != nil {
		files, err := ioutil.ReadDir("cmd")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				cmdName = file.Name()
			}
		}
		if _, err := os.Stat("cmd/" + cmdName); err != nil {
			log.Fatal("Unable to find entrypoint")
		}
	}

	if _, err := os.Stat(".env"); err == nil {
		bashenv.Source(".env")
	}

	if os.Getenv("GOPATH") == "" {
		log.Fatal("GOPATH must be set")
	}
	realCwd, _ := filepath.EvalSymlinks(cwd)
	if !strings.HasPrefix(realCwd, os.Getenv("GOPATH")) {
		log.Fatal("Must be run under GOPATH")
	}

	cmdPath := os.Getenv("GOPATH") + "/bin/" + cmdName
	cmd := NewCmdService(cmdPath, "-d", "dev/dev.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// build and run daemon first time
	buildAndRun(cmd, cmdName)

	startTime := time.Now()
	var runner *suture.Supervisor
	runner = suture.New("DevRunner", suture.Spec{
		Log: func(msg string) {
			if time.Now().Sub(startTime) < time.Duration(10*time.Second) {
				log.Println("Service failed before ready, exiting...")
				//webpackCmd.Stop()
				//automataCmd.Stop()
				os.Exit(1)
			}
			log.Println("Supervisor restarting a service...")
		},
	})

	runner.Add(cmd)

	if _, err := os.Stat("ui/package.json"); err == nil {
		webpackCmd := NewCmdService("npm", "run", "serve")
		webpackCmd.Dir = "ui"
		webpackCmd.Stdout = os.Stdout
		webpackCmd.Stderr = os.Stderr
		runner.Add(webpackCmd)
	}

	runner.ServeBackground()

	if os.Getenv("TIMEOUT") != "" {
		go func() {
			<-time.After(3 * time.Second)
			runner.Stop()
			os.Exit(0)
		}()
	}

	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		runner.Stop()
		os.Exit(0)
	}()

	// scan for go/config changes to rebuild and restart daemon
	notifyChanges("./...", []string{".go", ".toml"}, false, func(path string) {
		buildAndRun(cmd, cmdName)
	})
}

func buildAndRun(runner suture.Service, cmdName string) error {
	start := time.Now()
	cmd := exec.Command("go", "install", "./cmd/"+cmdName)
	output, err := cmd.CombinedOutput()
	if !cmd.ProcessState.Success() {
		buildErrored = true
		log.Println("ERROR! Build failed:")
		fmt.Println(string(output))
	} else {
		log.Println("New build:", time.Now().Sub(start))
		buildErrored = false
		runner.Stop()
	}
	time.Sleep(100 * time.Millisecond)
	return err
}

func extensionIn(path string, exts []string) bool {
	for _, ext := range exts {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

func notifyChanges(dir string, exts []string, onlyCreate bool, cb func(path string)) {
	c := make(chan notify.EventInfo, 1)
	types := notify.All
	if onlyCreate {
		types = notify.Create
	}
	if err := notify.Watch(dir, c, types); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	for event := range c {
		path := event.Path()
		dir, file := filepath.Split(path)
		if filepath.Base(dir) == ".git" {
			continue
		}
		if filepath.Base(file)[0] == '.' {
			continue
		}
		if extensionIn(path, exts) {
			cb(path)
		}
	}
}

type CmdService struct {
	*exec.Cmd
	signals chan os.Signal
	active  *exec.Cmd
}

func NewCmdService(name string, arg ...string) *CmdService {
	return &CmdService{
		Cmd: exec.Command(name, arg...),
	}
}

func (s *CmdService) Stop() {
	if s.active != nil && s.active.Process != nil {
		signal.Stop(s.signals)

		done := make(chan error)
		go func() {
			s.active.Wait()
			close(done)
			s.active = nil
		}()

		if runtime.GOOS == "windows" {
			s.active.Process.Kill()
		} else {
			s.active.Process.Signal(os.Interrupt)
		}

		select {
		case <-time.After(3 * time.Second):
			s.active.Process.Kill()
		case <-done:
			return
		}
		<-done
	}
}

func (s *CmdService) Serve() {
	// take a breather before starting a process
	time.Sleep(100 * time.Millisecond)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go func() {
		time.Sleep(100 * time.Millisecond)
		defer func() {
			signal.Stop(signals)
		}()
		for sig := range signals {
			if s.active == nil {
				return
			}
			s.active.Process.Signal(sig)
		}
	}()
	s.signals = signals
	cmd := &exec.Cmd{
		Path:   s.Path,
		Args:   s.Args,
		Env:    s.Env,
		Dir:    s.Dir,
		Stdin:  s.Stdin,
		Stdout: s.Stdout,
		Stderr: s.Stderr,
	}
	s.active = cmd
	s.active.Run()
}
