# zk-view
a `tree` like tool help you explore data structures in your zookeeper server

## Installation

```bash
go get github.com/dreamersdw/zk-view
go install github.com/dreamersdw/zk-view
```
# Usage
```
zk-view [--host=HOST] [--port=PORT] [--level=LEVEL] [--nodata] [--meta [--human]] [PATH]
zk-view --version
zk-view --help

Example:
	zk-view --host localhost /consumers`
```

## Screenshot

![zk-view](https://raw.githubusercontent.com/dreamersdw/zk-view/master/screenshot/zk-view.png)
