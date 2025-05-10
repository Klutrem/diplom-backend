package metrics

type NodeMetrics struct {
	NodeName              string   `json:"name"`
	CPUUsage              string   `json:"cpu_usage"`
	CpuCapacity           string   `json:"cpu_capacity"`
	CpuUsagePercentage    string   `json:"cpu_usage_percentage"`
	MemoryUsage           string   `json:"memory_usage"`
	MemoryCapacity        string   `json:"memory_capacity"`
	MemoryUsagePercentage string   `json:"memory_usage_percentage"`
	Roles                 []string `json:"roles"`
	Status                string   `json:"status"`
}

type PodMetrics struct {
	PodName     string `json:"pod_name"`
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
}
