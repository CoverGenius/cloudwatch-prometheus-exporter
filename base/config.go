package base

// Config represents the exporter configuration passed which is read at runtime from a YAML file.
type Config struct {
	Listen       string            `yaml:"listen,omitempty"`        // TCP Dial address for Prometheus HTTP API to listen on
	APIKey       string            `yaml:"api_key"`                 // AWS API Key ID
	APISecret    string            `yaml:"api_secret"`              // AWS API Secret
	Tags         []*TagDescription `yaml:"tags,omitempty"`          // Tags to filter resources by
	Period       uint8             `yaml:"period,omitempty"`        // How far back to request data for in minutes.
	Regions      []*string         `yaml:"regions"`                 // Which AWS regions to query resources and metrics for
	PollInterval uint8             `yaml:"poll_interval,omitempty"` // How often to fetch new data from the Cloudwatch API.
	LogLevel     uint8             `yaml:"log_level,omitempty"`     // Logging verbosity level
	Metrics      metric            `yaml:"metrics,omitempty"`       // Map of per metric configuration overrides
}

type metric struct {
	Data map[string]struct {
		Period int `yaml:"length,omitempty"` // How far back to request data for in minutes.
	} `yaml:",omitempty,inline"`
}
