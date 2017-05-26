package docker

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	cache "github.com/pmylund/go-cache"
)

func init() {
	com.Register("docker", &Component{},
		com.Option("name", "", "srv record for docker host discovery"),
		com.Option("listen", ":2475", "proxy listen address"),
		com.Option("version", client.DefaultVersion, "Docker client API version"))
}

const SessionHeaderKey = "X-Sandbox-Session"

type Component struct {
	sessCache *cache.Cache
	listener  *net.TCPListener
}

// Stop sandbox reverse proxy
func (c *Component) Stop() {
	// c.running = false
	if c.listener != nil {
		c.listener.Close()
	}
}

// Serve sandbox reverse proxy server
func (c *Component) Serve() {
	c.sessCache = cache.New(30*time.Second, 5*time.Minute)
	addr, _ := net.ResolveTCPAddr("tcp", com.GetString("listen"))
	var err error
	c.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err, addr.String())
	}
	defer c.listener.Close()
	log.Info("listening", log.Fields{
		"address": addr.String(),
	})

	for {
		conn, err := c.listener.AcceptTCP()
		if err != nil {
			log.Info(err)
			continue
		}
		go c.proxy(conn)
	}
}

func (c *Component) fetchBackend(id string) string {
	if v, found := c.sessCache.Get(id); id != "" && found {
		return v.(string)
	}
	return c.fetchNewBackend(id)
}

func (c *Component) fetchNewBackend(id string) string {
	_, addrs, err := net.LookupSRV("", "", com.GetString("name"))
	if err != nil {
		log.Info(err)
		return ""
	}
	target := strings.TrimSuffix(addrs[0].Target, ".")
	host := fmt.Sprintf("%s:%v", target, addrs[0].Port)
	c.sessCache.Set(id, host, cache.DefaultExpiration)
	return host
}

func (c *Component) proxy(conn *net.TCPConn) {
	defer conn.Close()
	reqID := strings.SplitN(uuid.New(), "-", 2)[0]
	fields := log.Fields{
		"requestID": reqID,
		"remote":    conn.RemoteAddr().String(),
	}

	r := bufio.NewReader(conn)
	req, err := http.ReadRequest(r)
	if err != nil {
		log.Info(fields, err)
		return
	}

	session := req.Header.Get(SessionHeaderKey)
	fields["sessionID"] = session
	host := c.fetchBackend(session)
	var backend net.Conn
	for attempt := 0; attempt < 5; attempt++ {
		log.Debug(fields, log.Fields{
			"backend": host,
			"attempt": strconv.Itoa(attempt),
		})
		backend, err = net.DialTimeout("tcp", host, 1*time.Second)
		if err != nil {
			log.Info(err, fields, log.Fields{
				"backend": host,
				"attempt": strconv.Itoa(attempt),
			})
			host = c.fetchNewBackend(session)
			continue
		}
		break
	}
	if err != nil {
		log.Info(fields, errors.Wrap(err, "failed to connect to a backend"))
		return
	}
	req.Write(backend)
	go func() {
		io.Copy(conn, backend)
		conn.CloseWrite()
		backend.(*net.TCPConn).CloseRead()
	}()
	io.Copy(backend, conn)
	backend.(*net.TCPConn).CloseWrite()
	conn.CloseRead()
	return
}
