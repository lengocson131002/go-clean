package bootstrap

type ServerConfig struct {
	Name               string
	AppVersion         string
	Port               int
	BaseURI            string
	GrRunningThreshold int // threshold for goroutines are running (which could indicate a resource leak).
	GcPauseThresholdMs int // threshold threshold garbage collection pause exceeds. (Millisecond)
	EnvFilePath        string
}
