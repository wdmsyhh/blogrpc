package log

type MetricLog struct {
	Type        string                 `json:"type"`
	Measurement string                 `json:"measurement"`
	Fields      map[string]uint64      `json:"fields"`
	Tags        map[string]interface{} `json:"tags"`
}

func NewMetricLog() *MetricLog {
	return &MetricLog{
		Type: "metric",
	}
}

func (log *MetricLog) ToJson() []byte {
	return ToJson(log)
}
