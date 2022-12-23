package log

type ServiceLog struct {
	Type      string                 `json:"type"`
	Level     string                 `json:"level"`
	Category  string                 `json:"category"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Backtrace string                 `json:"backtrace,omitempty"`
	ReqId     string                 `json:"reqId,omitempty"`
	Time      string                 `json:"time,omitempty"`
}

func NewServiceLog() *ServiceLog {
	return &ServiceLog{
		Type: "service",
	}
}

func (serviceLog *ServiceLog) ToJson() []byte {
	return ToJson(serviceLog)
}
