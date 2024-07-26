package main

import (
	"log"
	"os"
	"time"

	"blogrpc/practice/urfavecli/cmd"
	"blogrpc/practice/urfavecli/router"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalf("[APP] application run error: %s", err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "app"
	app.Version = "v1.2.10"
	app.Copyright = "(c) Gary"
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "gary.yin",
			Email: "",
		},
	}
	app.Writer = os.Stdout
	cli.ErrWriter = os.Stdout

	app.Commands = allCmd()
	return app
}

func allCmd() []*cli.Command {
	commands := cmd.Commands{
		Cmd1: cmd.NewCmd1(router.NewCmd1Router()),
		Cmd2: cmd.NewCmd2(router.NewCmd2Router()),
	}

	return []*cli.Command{
		commands.Cmd1,
		commands.Cmd2,
	}
}
