package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/mgutz/ansi"
	"github.com/samuel/go-zookeeper/zk"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	version = "0.1"
	usage   = `Usage:
	zk-view [--host=HOST] [--port=PORT] [--level=LEVEL] [--meta [--human]] [PATH]
	zk-view --version
	zk-view --help

Example:
	zk-view --host localhost /consumers`
)

var (
	zkHost      = "127.0.0.1"
	zkPort      = 2181
	zkPath      = "/"
	zkMeta      = false
	zkMaxLevel  = 1024
	zkHuman     = false
	turnOnColor = true
)

func colorize(s string, style string) string {
	if turnOnColor {
		return ansi.Color(s, style)
	}
	return s
}

func formatTimestamp(timestamp int64) interface{} {
	if zkHuman {
		return time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
	}
	return timestamp
}

func walk(root string, leading string, level int, conn *zk.Conn) {

	if level > zkMaxLevel {
		return
	}

	children, _, err := conn.Children(root)
	if err != nil {
		fmt.Printf("error, when get children of %s, %s\n", root, err)
		os.Exit(1)
	}

	childCount := len(children)
	for index, node := range children {
		isLast := index == childCount-1
		fullpath := path.Join(root, node)
		data, stat, _ := conn.Get(fullpath)
		displayNode(node, data, stat, leading, isLast)
		var extra string
		if isLast {
			extra = "    "
		} else {
			extra = "│   "
		}
		walk(fullpath, leading+extra, level+1, conn)
	}
}

func addGuideLine(multLine string, leading string) string {
	var result string
	lines := strings.Split(multLine, "\n")
	for index, line := range lines {
		if index != len(lines)-1 {
			result += leading + line + "\n"
		} else {
			result += leading + line
		}
	}
	return result
}

func displayNode(name string, data []byte, stat *zk.Stat, leading string, isLast bool) {
	var sep string
	if isLast {
		sep = "└── "
	} else {
		sep = "├── "
	}

	meta := map[string]interface{}{
		"Version":        stat.Version,
		"Cversion":       stat.Cversion,
		"Ctime":          formatTimestamp(stat.Ctime),
		"Mtime":          formatTimestamp(stat.Mtime),
		"EphemeralOwner": stat.EphemeralOwner,
		"DataLength":     stat.DataLength,
		"NumChildren":    stat.NumChildren,
	}

	var extra1 string
	if isLast {
		extra1 = "  "
	} else {
		extra1 = "│   "
	}

	var extra2, nodeColor string
	if stat.NumChildren > 0 {
		extra2 = "│"
		nodeColor = "blue"
	} else {
		extra2 = " "
		nodeColor = "green"
	}

	var metaData string
	if zkMeta {
		formatedMeta, _ := json.MarshalIndent(meta, leading+extra1+extra2, "  ")
		metaData = string(formatedMeta)
	} else {
		metaData = ""
	}

	fmt.Printf("%s%s%s %q %s\n", leading, sep, colorize(name, nodeColor), data, metaData)
}

func show(path string) {
	servers := []string{zkHost + ":" + strconv.FormatInt(int64(zkPort), 10)}
	conn, _, err := zk.Connect(servers, 100*time.Second)
	if err != nil {
		fmt.Printf("meet an error when connect to %s\n ", servers)
	}

	walk(path, "", 1, conn)
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

	if opt["--level"] != nil {
		level, _ := strconv.ParseInt(opt["--level"].(string), 10, 32)
		zkMaxLevel = int(level)
	}

	zkMeta = opt["--meta"].(bool)
	zkHuman = opt["--human"].(bool)
	turnOnColor = terminal.IsTerminal(int(os.Stdout.Fd()))

	fmt.Println(zkPath)
	show(zkPath)
}
