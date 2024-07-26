//go:build wireinject
// +build wireinject

package main

import (
	"blogrpc/practice/urfavecli/cmd"
	"blogrpc/practice/urfavecli/router"
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

var providerSet = wire.NewSet(
	// router
	router.NewCmd1Router,
	router.NewCmd2Router,

	// cmd
	cmd.ProviderSet,
)

func initApp() (cli.Commands, func()) {
	panic(wire.Build(providerSet))
}
