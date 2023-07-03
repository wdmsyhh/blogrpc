## 注意

- 命令行参数解析spf13/pflag包 | 标准库 flag包:https://blog.csdn.net/inthat/article/details/123623603#:~:text=pflag%E6%98%AFGo%E7%9A%84%E6%9C%AC%E6%9C%BA%E6%A0%87%E5%BF%97%E5%8C%85%E7%9A%84%E7%9B%B4%E6%8E%A5%E6%9B%BF%E4%BB%A3%E3%80%82%20%E5%A6%82%E6%9E%9C%E6%82%A8%E5%9C%A8%E5%90%8D%E7%A7%B0%E2%80%9C%20flag%E2%80%9D%E4%B8%8B%E5%AF%BC%E5%85%A5pflag%20%28%E5%A6%82%EF%BC%9A%20import%20flag,%22github.com%2Fspf13%2Fpflag%22%29%EF%BC%8C%E5%88%99%E6%89%80%E6%9C%89%E4%BB%A3%E7%A0%81%E5%BA%94%E7%BB%A7%E7%BB%AD%E8%BF%90%E8%A1%8C%E4%B8%94%E6%97%A0%E9%9C%80%E6%9B%B4%E6%94%B9%E3%80%82%20%E4%B8%80%E4%B8%AA%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E5%9C%A8%20Pflag%20%E5%8C%85%E4%B8%AD%E4%BC%9A%E8%A2%AB%E8%A7%A3%E6%9E%90%E4%B8%BA%E4%B8%80%E4%B8%AA%20Flag%20%E7%B1%BB%E5%9E%8B%E7%9A%84%E5%8F%98%E9%87%8F%E3%80%82

- cobra 和 标准 flag 使用的时候好像有点问题

如果使用标准的 "flag" 包，下面代码执行 `go run main.go --job=1 --env=p cmd1 cmd11 b c` 的时候会打印如下：
```
product init
[/tmp/go-build1749531715/b001/exe/main --job=1 --env=p cmd1 cmd11 b c]
p
true
Error: unknown flag: --job
Usage:
  rootCmd cmd1 cmd11 [flags]

Flags:
  -h, --help   help for cmd11

```

使用 `flag "github.com/spf13/pflag"` 的时候打印正常：

```
product init
[/tmp/go-build1983867991/b001/exe/main --job=1 --env=p cmd1 cmd11 b c]
p
true
exec cmd11, args: [b c]

```

使用 `go run main.go --job=1 --env=p   cmd2 b c d` 的时候执行的是 cmd2 命令:
```
product init
[/tmp/go-build1752896154/b001/exe/main --job=1 --env=p cmd2 b c c]
p
true
exec cmd2, args: [b c d]
```

```go
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
```