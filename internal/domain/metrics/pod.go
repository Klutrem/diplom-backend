package metrics

type PodMetrics struct {
	PodName            string   `json:"pod_name"`
	Namespace          string   `json:"namespace"`
	NodeName           string   `json:"node_name"`
	CPUUsage           int64    `json:"cpu_usage"`
	CPUUsagePercent    *float64 `json:"cpu_usage_percent"`
	CPUUsageLimit      int64    `json:"cpu_usage_limit"`
	CPUUsageRequest    int64    `json:"cpu_usage_request"`
	MemoryUsage        int64    `json:"memory_usage"`
	MemoryUsagePercent *float64 `json:"memory_usage_percent,omitempty"`
	MemoryUsageLimit   int64    `json:"memory_usage_limit"`
	MemoryUsageRequest int64    `json:"memory_usage_request"`
	Status             string   `json:"status"`
	StartTime          string   `json:"start_time"`
	RestartCount       int32    `json:"restart_count"`
}
