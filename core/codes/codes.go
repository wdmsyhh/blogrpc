package codes

type Code uint32

// The canonical error codes used by core gRPC service.
const (
	SystemError         Code = 1000001 // system unexpected exception
	ServiceUnavailable  Code = 1000002 // service is not available
	RemoteServiceError  Code = 1000003 // remote service error occurs
	InternalServerError Code = 1000004
	UnknownError        Code = 1000005
	InvalidParams       Code = 1000006
)

const (
	SEPARATOR string = "\n\n"
)
