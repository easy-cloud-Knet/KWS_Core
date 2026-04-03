package status

type SourceType string

const (
	CPU       SourceType = "cpu"
	Memory    SourceType = "memory"
	MaxMemory SourceType = "max_memory"
	CPUTime   SourceType = "cpu_time"
)
