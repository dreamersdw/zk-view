package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/mgutz/ansi"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	version = "0.1"
	usage   = `Usage:
	zk-view [--host=HOST] [--port=PORT] [PATH]
	zk-view --version
	zk-view --help

Example:
	zk-view 'tasks:*' 'metrics:*' `
)

var (
	zkHost      = "127.0.0.1"
	zkPort      = 2181
	zkPath      = "/"
	turnOnColor = true
)

func colorize(s string, style string) string {
	if turnOnColor {
		return ansi.Color(s, style)
	}
	return s
}

func walk(root string, leading string, conn *zk.Conn) {
	children, _, err := conn.Children(root)
	if err != nil {
		fmt.Printf("error, when get children of %s, %s\n", root, err)
		os.Exit(1)
	}

	childCount := len(children)
	for index, node := range children {
		isLast := index == childCount-1
		var extra string

		if isLast {
			extra = "    "
		} else {
			extra = "│   "
		}

		fullpath := path.Join(root, node)

		data, stat, _ := conn.Get(fullpath)
		displayNode(node, data, stat, leading, isLast)
		walk(fullpath, leading+extra, conn)
	}
}

func displayNode(name string, data []byte, stat *zk.Stat, leading string, isLast bool) {
	var sep string
	if isLast {
		sep = "└── "
	} else {
		sep = "├── "
	}
	fmt.Printf("%s%s%s # %s\n", leading, sep, colorize(name, "blue"), string(data))
}

func show(path string) {
	servers := []string{zkHost + ":" + strconv.FormatInt(int64(zkPort), 10)}
	conn, _, err := zk.Connect(servers, 10*time.Second)
	if err != nil {
		fmt.Printf("meet an error when connect to %s\n", servers)
	}

	walk(path, "", conn)
}

func main() {
	opt, err := docopt.Parse(usage, nil, false, "", false, false)
	if err != nil {
		fmt.Println("error when parse cmdline")
		os.Exit(1)
	}

	if opt["PATH"] != nil {
		zkPath = opt["PATH"].(string)
	}

	if opt["--host"] != nil {
		zkHost = opt["--host"].(string)
	}

	if opt["--port"] != nil {
		port, _ := strconv.ParseInt(opt["--port"].(string), 10, 32)
		zkPort = int(port)
	}

	fmt.Println(zkPath)
	show(zkPath)
}
