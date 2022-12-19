--go_out参数用来指定 protoc-gen-go 插件的工作方式和go代码的生成位置

--go_out主要的两个参数为 plugins 和 paths，分别表示生成go代码所使用的插件和生成的go代码的位置。

--go_out的写法是参数之间用 逗号 隔开，最后加上 冒号 来指定代码的生成位置。

比如：--go_out=plugins=grpc,paths=import:.

paths参数有两个选项，分别是 import 和 source_relative， 默认为 import ，表示按照生成的go代码的包的全路径去创建目录层级，source_relative 表示按照proto源文件的目录层级去创建go代码的目录层级。