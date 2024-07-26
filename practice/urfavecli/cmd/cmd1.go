package cmd

import (
	"fmt"

	"blogrpc/practice/urfavecli/router"
	"github.com/urfave/cli/v2"
)

type Cmd1 *cli.Command

func NewCmd1(router *router.Cmd1Router) Cmd1 {
	return &cli.Command{
		Name:  "cmd1",
		Usage: "cmd1 command eg: ./app cmd1 --addr=:8081",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "addr",
				Usage:    "--addr=:8081",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "--name=你好$小明",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			addr := c.String("addr")

			name := c.String("name")
			fmt.Println(name)

			return router.Run(c.Context, addr)
		},
	}
}
