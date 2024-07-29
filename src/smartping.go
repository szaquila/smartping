// 用去CGO的驱动  CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/smartping src/smartping.go
// centos 6.1     go build -tags libsqlite3 -ldflags '-s -w' -o bin/smartping src/smartping.go
// alpine(docker) CC=musl-gcc go build -tags libsqlite3 -tags musl -ldflags '-linkmode external -extldflags "-static" -s -w' -o bin/smartping src/smartping.go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/jakecoffman/cron"
	"github.com/smartping/smartping/src/funcs"
	"github.com/smartping/smartping/src/g"
	"github.com/smartping/smartping/src/http"
	//"sync"
)

// Init config
var Version = "0.8.0"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	g.ParseConfig(Version)
	go funcs.ClearArchive()
	c := cron.New()
	c.AddFunc("*/60 * * * * *", func() {
		go funcs.Ping()
		go funcs.Mapping()
		if g.Cfg.Mode["Type"] == "cloud" {
			go funcs.StartCloudMonitor()
		}
	}, "ping")
	c.AddFunc("0 0 * * * *", func() {
		go funcs.ClearArchive()
	}, "mtc")
	c.Start()
	http.StartHttp()
}
