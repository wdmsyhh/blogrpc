package main

import (
	"fmt"
	"github.com/spf13/cobra"
	// "flag"
	flag "github.com/spf13/pflag"
	"os"
)

var (
	env = flag.String("env", "local", "the running environment")
	// job flags
	enableJob = flag.Bool("job", false, "blogrpc job name")
)

func init() {
	RootCmd.AddCommand(cmd1)
	RootCmd.AddCommand(cmd2)
	cmd1.AddCommand(cmd11)
}

func main() {
	flag.Parse() // 重要，重要，重要
	fmt.Println("product init")
	fmt.Println(os.Args)
	fmt.Println(*env)
	fmt.Println(*enableJob)
	RootCmd.Execute()
	select {}
}

var RootCmd = &cobra.Command{
	Use:   "rootCmd",
	Short: "I am root",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("root cmd, args:", args)
		return nil
	},
}

var cmd1 = &cobra.Command{
	Use:   "cmd1",
	Short: "I am cmd1",
	Long:  `I am cmd1......`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("exec cmd1, args:", args)
		return nil
	},
}

var cmd11 = &cobra.Command{
	Use:   "cmd11",
	Short: "I am cmd11",
	Long:  `I am cmd11......`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("exec cmd11, args:", args)
		return nil
	},
}

var cmd2 = &cobra.Command{
	Use:   "cmd2",
	Short: "I am cmd2",
	Long:  `I am cmd2......`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("exec cmd2, args:", args)
		return nil
	},
}
