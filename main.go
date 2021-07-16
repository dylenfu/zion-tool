package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dylenfu/zion-tool/journal"
	"github.com/dylenfu/zion-tool/utils/files"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"
)

var (
	CmdTPS = cli.Command{
		Name:   "tps",
		Usage:  "try to test zion consensus tps.",
		Action: journal.HandleTPS,
	}
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "hotstuff test tool"
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2021 The Ontology Authors"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		CmdTPS,
	}
	app.Before = beforeCommands
	app.After = afterCommands
	return app
}

func main() {

	app := setupApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// action execute after commands
func beforeCommands(ctx *cli.Context) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return nil
}

func afterCommands(ctx *cli.Context) error {
	log.Info("\r\n" +
		"\r\n" +
		"\r\n")
	return nil
}
