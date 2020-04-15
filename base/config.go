package base

import (
	"strings"

	"github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type configMetric struct {
	AWSMetric     string                  `yaml:"metric"`         // The Cloudwatch metric to use
	Help          string                  `yaml:"help"`           // Custom help text for the generated metric
	Dimensions    []*cloudwatch.Dimension `yaml:"dimensions"`     // The resource dimensions to generate individual series for (via labels)
	Statistics    []*string               `yaml:"statistics"`     // List of AWS statistics to use.
	OutputName    string                  `yaml:"output_name"`    // Allows override of the generate metric name
	RangeSeconds  int64                   `yaml:"range_seconds"`  // How far back to request data for in seconds.
	PeriodSeconds int64                   `yaml:"period_seconds"` // Granularity of results from cloudwatch API.
}

type metric struct {
	Data map[string][]*configMetric `yaml:",omitempty,inline"` // Map from namespace to list of metrics to scrape.
}

// Config represents the exporter configuration passed which is read at runtime from a YAML file.
type Config struct {
	Listen       string            `yaml:"listen,omitempty"`        // TCP Dial address for Prometheus HTTP API to listen on
	APIKey       string            `yaml:"api_key"`                 // AWS API Key ID
	APISecret    string            `yaml:"api_secret"`              // AWS API Secret
	Tags         []*TagDescription `yaml:"tags,omitempty"`          // Tags to filter resources by
	Regions      []*string         `yaml:"regions"`                 // Which AWS regions to query resources and metrics for
	LogLevel     uint8             `yaml:"log_level,omitempty"`     // Logging verbosity level
	PollInterval int64             `yaml:"poll_interval,omitempty"` // How often to fetch new data from the Cloudwatch API.

	// Default values for metrics, will only be used for a metric if that
	// metric does not have an override configured
	PeriodSeconds int64 `yaml:"period_seconds,omitempty"` // Granularity of results from cloudwatch API.
	RangeSeconds  int64 `yaml:"range_seconds,omitempty"`  // How far back to request data for.

	Metrics metric `yaml:"metrics"` // Map of per metric configuration overrides
}

// ConstructMetrics generates a map of MetricDescriptions keyed by CloudWatch namespace using the defaults provided in Config.
func (c *Config) ConstructMetrics(defaults map[string]map[string]*MetricDescription) map[string][]*MetricDescription {
	mds := make(map[string][]*MetricDescription)
	for namespace, metrics := range c.Metrics.Data {

		if len(metrics) <= 0 {
			if namespaceDefaults, ok := defaults[namespace]; ok {
				for key, defaultMetric := range namespaceDefaults {
					metrics = append(metrics, &configMetric{
						AWSMetric:     key,
						OutputName:    *defaultMetric.OutputName,
						Help:          *defaultMetric.Help,
						PeriodSeconds: defaultMetric.PeriodSeconds,
						RangeSeconds:  defaultMetric.RangeSeconds,
						Dimensions:    defaultMetric.Dimensions,
					})
				}
			}
		}

		mds[namespace] = []*MetricDescription{}
		for _, metric := range metrics {
			name := metric.OutputName
			if name == "" {
				name = helpers.ToPromString(strings.TrimPrefix(namespace, "AWS/") + "_" + metric.AWSMetric)
			}

			period := metric.PeriodSeconds
			if period == 0 {
				period = c.PeriodSeconds
			}

			rangeSeconds := metric.RangeSeconds
			if rangeSeconds == 0 {
				rangeSeconds = c.RangeSeconds
			}

			if metric.Statistics == nil || len(metric.Statistics) < 1 {
				metric.Statistics = helpers.StringPointers("Average")
			}

			help := metric.Help
			if help == "" {
				if d, ok := defaults[namespace][metric.AWSMetric]; ok {
					help = *d.Help
				}
			}

			// TODO handle dimensions
			// TODO move metricName function here / apply to output name
			// TODO one stat per metric
			mds[namespace] = append(mds[namespace], &MetricDescription{
				Help:          &help,
				OutputName:    &name,
				Dimensions:    metric.Dimensions,
				PeriodSeconds: period,
				RangeSeconds:  rangeSeconds,
				Statistic:     metric.Statistics,

				Namespace: namespace,
				AWSMetric: metric.AWSMetric,
			})
		}
	}
	return mds
}
