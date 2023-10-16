package constant

import "blogrpc/core/util"

var (
	// for in docker test
	SERVICE_HELLO_PORT  = "1701"
	SERVICE_MEMBER_PORT = "1701"

	SERVICE_HELLO_HOST  = "blogrpc-hello"
	SERVICE_MEMBER_HOST = "blogrpc-member"

	// for local test
	//SERVICE_HELLO_PORT  = "1701"
	//SERVICE_MEMBER_PORT = "1702"
	//
	//SERVICE_HELLO_HOST  = "localhost"
	//SERVICE_MEMBER_HOST = "localhost"
)

func init() {
	if !util.IsRunningInContainer() {
		// for local test
		SERVICE_HELLO_PORT = "1701"
		SERVICE_MEMBER_PORT = "1702"

		SERVICE_HELLO_HOST = "localhost"
		SERVICE_MEMBER_HOST = "localhost"
	}
}
