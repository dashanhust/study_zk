package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-zookeeper/zk"
)

const flag int32 = 0

var acl = zk.WorldACL(zk.PermAll)
var ErrUnknownCmd = errors.New("unknown command")

type Cmd struct {
	Name        string
	Options     []string
	ExitWhenErr bool
	Conn        *zk.Conn
	Config      *Config
}

func NewCmd(name string, options []string, conn *zk.Conn, config *Config) *Cmd {
	return &Cmd{
		Name:    name,
		Options: options,
		Conn:    conn,
		Config:  config,
	}
}

func (c *Cmd) connected() bool {
	state := c.Conn.State()
	return state == zk.StateConnected || state == zk.StateHasSession
}

func (c *Cmd) checkConn() (err error) {
	if !c.connected() {
		err = errors.New("connection is disconnected")
	}
	return
}

func ParseCmd(cmd string) (name string, options []string) {
	args := make([]string, 0)
	for _, c := range strings.Split(cmd, " ") {
		if c != "" {
			args = append(args, c)
		}
	}
	if len(args) != 0 {
		return
	}
	return args[0], args[1:]
}

// Cmd: ls 列出指定路径的节点信息
func (c *Cmd) ls() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	children, _, err := c.Conn.Children(p)
	if err != nil {
		return
	}
	fmt.Printf("[%s]\n", strings.Join(children, ", "))
	return
}

// Cmd: get 获取指定节点的内容
func (c *Cmd) get() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	value, stat, err := c.Conn.Get(p)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n%s\n", string(value), fmtStat(stat))
	return
}

// Cmd: create 创建节点
func (c *Cmd) create() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	p = cleanPath(p)
	_, err = c.Conn.Create(p, []byte(data), flag, acl)
	if err != nil {
		return
	}
	fmt.Printf("Created %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

// Cmd: set 设置节点的值
func (c *Cmd) set() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	p = cleanPath(p)
	stat, err := c.Conn.Set(p, []byte(data), -1)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", fmtStat(stat))
	return
}

// Cmd: delete 指定的节点
func (c *Cmd) delete() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	err = c.Conn.Delete(p, -1)
	if err != nil {
		return
	}
	fmt.Printf("Deleted %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

// Cmd: close 关闭zk链接
func (c *Cmd) close() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	c.Conn.Close()
	if !c.connected() {
		fmt.Println("Closed")
	}
	return
}

// Cmd: connect 创建zk链接
func (c *Cmd) connect() (err error) {
	options := c.Options
	var conn *zk.Conn
	if len(options) > 0 {
		cf := NewConfig(strings.Split(options[0], ","), false)
		conn, err = cf.Connect()
	} else {
		conn, err = c.Config.Connect()
	}
	if err != nil {
		return
	}

	if c.connected() {
		c.Conn.Close()
	}
	c.Conn = conn
	fmt.Println("Connected")
	return
}

// Cmd: addAuth 添加权限
func (c *Cmd) addAuth() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	options := c.Options
	if len(options) < 2 {
		return errors.New("addAuth <scheme> <auth>")
	}
	scheme := options[0]
	auth := options[1]
	err = c.Conn.AddAuth(scheme, []byte(auth))
	if err != nil {
		return
	}
	fmt.Println("Added")
	return
}

// Cmd: run 执行命令
func (c *Cmd) run() (err error) {
	switch c.Name {
	case "ls":
		return c.ls()
	case "get":
		return c.get()
	case "create":
		return c.create()
	case "set":
		return c.set()
	case "delete":
		return c.delete()
	case "close":
		return c.close()
	case "connect":
		return c.connect()
	case "addauth":
		return c.addAuth()
	default:
		return ErrUnknownCmd
	}
}

// Cmd: Run 供外部调用
func (c *Cmd) Run() {
	err := c.run()
	if err != nil {
		if err == ErrUnknownCmd {
			printHelp()
			if c.ExitWhenErr {
				os.Exit(2)
			}
		} else {
			printRunError(err)
			if c.ExitWhenErr {
				os.Exit(3)
			}
		}
	}
}

func GetExecutor(cmd *Cmd) func(s string) {
	return func(s string) {
		name, options := ParseCmd(s)
		cmd.Name = name
		cmd.Options = options
		if name == "exit" {
			os.Exit(0)
		}
		cmd.Run()
	}
}
