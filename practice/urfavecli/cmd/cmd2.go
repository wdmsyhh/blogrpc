package cmd

import (
	"blogrpc/practice/urfavecli/router"
	"github.com/urfave/cli/v2"
)

type Cmd2 *cli.Command

func NewCmd2(router *router.Cmd2Router) Cmd2 {
	return &cli.Command{
		Name:  "cmd2",
		Usage: "cmd2 command eg: ./app cmd2 --addr=:8082",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "addr",
				Usage:    "--addr=:8082",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			addr := c.String("addr")

			return router.Run(c.Context, addr)
		},
	}
}
