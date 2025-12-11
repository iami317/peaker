package main

import (
	"github.com/iami317/logx"
	"github.com/iami317/peaker"
	"github.com/iami317/peaker/plugins"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func main() {
	RunApp()
}

func RunApp() {
	app := cli.NewApp()
	app.Usage = ""
	app.Name = "Peaker"
	app.Version = "1.0 beta"
	app.Description = ""
	app.HelpName = "./peaker -h"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "ip_list",
			Aliases: []string{"i"},
			Value:   "./iplist.txt",
		},
		&cli.StringFlag{
			Name:    "user_dict",
			Aliases: []string{"u"},
			Value:   "./user.dic",
		},
		&cli.StringFlag{
			Name:    "pass_dict",
			Aliases: []string{"p"},
			Value:   "./pass.dic",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Value: true,
			Usage: "Set log level to debug",
		},
		&cli.BoolFlag{
			Name:    "check_alive",
			Aliases: []string{"cA"},
			Usage:   "Check if the target is alive",
		},
		&cli.IntFlag{
			Name:    "thread",
			Aliases: []string{"c"},
			Value:   30,
			Usage:   "Number of concurrent threads",
		},
		&cli.IntFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   50000,
			Usage:   "The maximum execution time of a single ip",
		},
		&cli.IntFlag{
			Name:    "timeout-single",
			Aliases: []string{"tS"},
			Value:   30,
			Usage:   "The maximum execution time of a single account",
		},
		&cli.IntFlag{
			Name:    "thread-single",
			Aliases: []string{"tC"},
			Value:   100,
			Usage:   "The number of concurrency for a single protocol",
		},
		&cli.BoolFlag{
			Name:  "protocol",
			Usage: "view supported protocols",
		},
	}
	app.Action = RunServer
	err := app.Run(os.Args)
	if err != nil {
		logx.Fatalf("engin err: %v", err)
		return
	}
}

func RunServer(ctx *cli.Context) error {
	if ctx.Bool("protocol") {
		for protocol, _ := range plugins.ScanMap {
			logx.Silent(string(protocol))
		}
		return nil
	}
	config := peaker.Config{
		Logger: logx.New(),
	}
	if ctx.IsSet("verbose") {
		config.Logger.SetLevel("verbose")
		config.DebugMode = true
	}

	if ctx.IsSet("check_alive") {
		config.CheckAlive = true
	}

	if ctx.IsSet("thread-single") {
		config.ThreadSingle = ctx.Int("thread-single")
	}
	config.Thread = ctx.Int("thread")
	config.TimeOut = time.Duration(ctx.Int("timeout")) * time.Second
	config.Ts = time.Duration(ctx.Int("tS")) * time.Second

	w := peaker.NewWeak(config)
	w.StartTime = time.Now()

	userDict, err := w.ReadUserDict(ctx.String("user_dict"))
	if err != nil {
		return err
	}
	passDict, err := w.ReadPasswordDict(ctx.String("pass_dict"))
	if err != nil {
		return err
	}
	ipList, err := w.ReadIpList(ctx.String("ip_list"))
	if err != nil {
		return err
	}
	resultChan := make(chan interface{}, 1)
	w.RunTask(ipList, userDict, passDict, resultChan)
	for v := range resultChan {
		r := v.(*peaker.ResultOut)
		if len(r.Crack) > 0 {
			for _, crack := range r.Crack {
				logx.Silent(r.Addr.String(), crack.String())
			}
		}
	}
	return nil
}
