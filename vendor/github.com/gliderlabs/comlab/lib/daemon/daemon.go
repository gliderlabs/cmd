package daemon

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/com/viper"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/mgutz/ansi"
	"github.com/spf13/cast"
	"github.com/thejerf/suture"
)

func init() {
	if !LocalMode() {
		ansi.DisableColors(true)
	}
}

var (
	gray   = ansi.ColorFunc("black+h")
	cyan   = ansi.ColorFunc("cyan")
	red    = ansi.ColorFunc("red")
	yellow = ansi.ColorFunc("yellow")
	reset  = ansi.ColorFunc("reset")
	bright = ansi.ColorFunc("white+h")
)

type consoleOutput struct{}

func (c *consoleOutput) Log(e log.Event) {
	color := reset
	switch e.Type {
	case log.TypeLocal:
		color = yellow
	case log.TypeFatal:
		color = red
	case log.TypeInfo:
		if DebugMode() {
			color = bright
		}
	}
	if _, ok := e.Fields["err"]; ok {
		color = red
	}
	pkg := e.Fields["pkg"]
	e = e.Remove("pkg")
	var parts []string
	for _, key := range e.Index {
		if key == "msg" {
			parts = append([]string{e.Fields[key]}, parts...)
		} else {
			parts = append(parts, fmt.Sprintf("%s=%v", key, e.Fields[key]))
		}
	}
	fmt.Println(gray(e.Time.Format("15:04:05.000")), cyan("["+pkg+"]"), color(strings.Join(parts, " ")))
}

type LifecycleContributor interface {
	AppPreStart() error
}

func LocalMode() bool {
	return os.Getenv("LOCAL") != "false"
}

func DebugMode() bool {
	return os.Getenv("DEBUG") != ""
}

func Run(name string) {
	log.RegisterObserver(new(consoleOutput))
	log.SetFieldProcessor(func(e log.Event, o interface{}) (log.Event, bool) {
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
		}
		return e, false
	})

	cfg := viper.NewConfig()
	com.SetConfig(cfg)
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-d" {
				os.Setenv("DEBUG", "1")
			}
			if !strings.HasPrefix(arg, "-") {
				cfg.SetConfigFile(arg)
				err := cfg.ReadInConfig()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.SetDebug(DebugMode())
	log.SetLocal(LocalMode())

	for _, service := range com.Enabled(new(LifecycleContributor), nil) {
		if err := service.(LifecycleContributor).AppPreStart(); err != nil {
			log.Fatal(err)
		}
	}

	app := suture.NewSimple(name)
	for _, service := range com.Enabled(new(suture.Service), nil) {
		app.Add(service.(suture.Service))
	}
	app.Serve()
}
