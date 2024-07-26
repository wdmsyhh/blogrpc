package cmd

import (
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

type Commands struct {
	Cmd1 Cmd1
	Cmd2 Cmd2
}

// NewCommands 创建cli命令
func NewCommands(cmd *Commands) cli.Commands {
	return cli.Commands{
		cmd.Cmd1,
		cmd.Cmd2,
	}
}

// ProviderSet 用于生成这样的结构：
//
//	Commands{
//		Cmd1: NewCmd1(cmd1Router),
//		Cmd2: NewCmd2(cmd2Router),
//	}
var ProviderSet = wire.NewSet(
	NewCmd1,
	NewCmd2,
	NewCommands,
	wire.Struct(new(Commands), "*"),
)
