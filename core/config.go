package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type Auth struct {
	Scheme  string
	Payload []byte
}

func NewAuth(scheme string, auth string) *Auth {
	return &Auth{
		Scheme:  scheme,
		Payload: []byte(auth),
	}
}

type Config struct {
	Servers   []string
	Auth      *Auth
	SilentLog bool
}

func NewConfig(servers []string, silentLog bool) *Config {
	return &Config{
		Servers:   servers,
		SilentLog: silentLog,
	}
}

type emptyLogger struct {
}

func (emptyLogger) Printf(format string, a ...interface{}) {
}

func (c *Config) Connect() (conn *zk.Conn, err error) {
	logger := zk.WithLogger(zk.DefaultLogger)
	if c.SilentLog {
		logger = zk.WithLogger(emptyLogger{})
	}
	conn, e, err := zk.Connect(c.Servers, time.Second, logger)
	if err != nil {
		return
	}
	if c.Auth != nil {
		auth := c.Auth
		err = conn.AddAuth(auth.Scheme, auth.Payload)
		if err != nil {
			return
		}
	}

	n := 0
	failed := false
loop:
	for {
		event, ok := <-e
		n += 1
		if ok && event.State == zk.StateConnected {
			break loop
		} else if n > 3 {
			failed = true
			break loop
		}
	}

	if failed {
		err = fmt.Errorf("failed to connect to %s", strings.Join(c.Servers, ","))
	}
	return
}
