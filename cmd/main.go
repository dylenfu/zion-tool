package main

import (
	"flag"
	"math/rand"
	"strings"
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/core"
	"github.com/dylenfu/zion-tool/pkg/frame"
	"github.com/dylenfu/zion-tool/pkg/log"
)

var (
	loglevel   int    // log level [1: debug, 2: info]
	configpath string //config file
	Methods    string //Methods list in cmdline
)

func init() {
	flag.StringVar(&configpath, "config", "config.json", "configpath of palette-tool")
	flag.StringVar(&Methods, "t", "connect", "methods to run. use ',' to split methods")
	flag.IntVar(&loglevel, "loglevel", 2, "loglevel [1: debug, 2: info]")

	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	defer time.Sleep(time.Second)

	log.InitLog(loglevel, log.Stdout)
	config.Init(configpath)
	core.Endpoint()

	methods := make([]string, 0)
	if Methods != "" {
		methods = strings.Split(Methods, ",")
	}

	frame.Tool.Start(methods)
}
