package constant

import "blogrpc/core/util"

// 此处定义全局变量为以容器运行时使用端口和服务名
var (
	ServiceHelloPort = "1701"
	ServiceHelloHost = "blogrpc-hello"

	ServiceMemberPort = "1701"
	ServiceMemberHost = "blogrpc-member"
)

// 本地直接运行时初始化全局变量为本地使用的端口和服务名
func init() {
	// 本地调试是修改host和port
	if !util.IsRunningInContainer() {
		// hello 服务
		ServiceHelloPort = "1701"
		ServiceHelloHost = "localhost"
		// member 服务
		ServiceMemberPort = "1702"
		ServiceMemberHost = "localhost"
	}
}
