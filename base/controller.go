package base

import (
	"fmt"
	"strings"
	"sync"
	"time"

	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	results    = make(map[string]prometheus.Collector)
	timestamps = make(map[prometheus.Collector]*time.Time)
)

type TimeDescription struct {
	StartTime *time.Time
	EndTime   *time.Time
}

type TagDescription struct {
	Key   *string `yaml:"name"`
	Value *string `yaml:"value"`
}

type DimensionDescription struct {
	Name  *string
	Value *string
}

type metric struct {
	Data map[string]struct {
		Length int `yaml:"length,omitempty"`
	} `yaml:",omitempty,inline"`
}

type Config struct {
	Listen       string            `yaml:"listen,omitempty"`
	APIKey       string            `yaml:"api_key"`
	APISecret    string            `yaml:"api_secret"`
	Tags         []*TagDescription `yaml:"tags,omitempty"`
	Period       uint8             `yaml:"period,omitempty"`
	Regions      []*string         `yaml:"regions"`
	PollInterval uint8             `yaml:"poll_interval,omitempty"`
	LogLevel     uint8             `yaml:"log_level,omitempty"`
	Metrics      metric            `yaml:"metrics,omitempty"`
}

type Metrics struct {
	Mutex  sync.RWMutex
	Metric map[string]map[string]*MetricDescription
}

type MetricDescription struct {
	Help       *string
	Type       *string
	OutputName *string
	Dimensions []*cloudwatch.Dimension
	Period     int
	Statistic  []*string
	Data       map[string][]*string
}

type RegionDescription struct {
	Config     *Config
	Session    *session.Session
	Tags       []*TagDescription
	Region     *string
	AccountID  *string
	Filters    []*ec2.Filter
	Namespaces map[string]*NamespaceDescription
	Time       *TimeDescription
	Mutex      sync.RWMutex
	Period     *uint8
}

type NamespaceDescription struct {
	Namespace *string
	Resources []*ResourceDescription
	Parent    *RegionDescription
	Mutex     sync.RWMutex
	Metrics   map[string]*MetricDescription
}

type ResourceDescription struct {
	Name       *string
	ID         *string
	Dimensions []*cloudwatch.Dimension
	Type       *string
	Parent     *NamespaceDescription
	Mutex      sync.RWMutex
	Query      []*cloudwatch.MetricDataQuery
}

func (md *MetricDescription) metricName(stat string) *string {
	suffix := ""
	switch stat {
	case "Average":
		// For backwards compatibility we have to omit the _avg
		suffix = ""
	case "Sum":
		suffix = "_sum"
	case "Minimum":
		suffix = "_min"
	case "Maximum":
		suffix = "_max"
	case "SampleCount":
		suffix = "_count"
	}
	name := *md.OutputName + suffix
	return &name
}

func (rd *RegionDescription) SetRequestTime() error {
	log.Debug("Setting request time ...")
	td := TimeDescription{}
	t := time.Now().Round(time.Minute * 5)
	start := t.Add(time.Minute * -time.Duration(*rd.Period))
	td.StartTime = &start
	td.EndTime = &t
	rd.Time = &td
	return nil
}

func (rd *RegionDescription) BuildArn(s *string, r *string) (string, error) {
	a := arn.ARN{
		Service:   *s,
		Region:    *rd.Region,
		AccountID: *rd.AccountID,
		Resource:  *r,
		Partition: "aws",
	}
	return a.String(), nil
}

func (rd *RegionDescription) BuildFilters() error {
	filters := []*ec2.Filter{}
	for _, tag := range rd.Tags {
		f := &ec2.Filter{
			Name:   aws.String(strings.Join([]string{"tag", *tag.Key}, ":")),
			Values: []*string{tag.Value},
		}
		filters = append(filters, f)
	}
	rd.Filters = filters
	return nil
}

func (rd *RegionDescription) GetAccountId() error {
	session := iam.New(rd.Session)
	input := iam.GetUserInput{}
	user, err := session.GetUser(&input)
	h.LogError(err)
	a, err := arn.Parse(*user.User.Arn)
	h.LogError(err)
	rd.AccountID = &a.AccountID

	return nil
}

func (rd *RegionDescription) Init(s *session.Session, td []*TagDescription, r *string, p *uint8) error {
	log.Infof("Initializing region %s ...", *r)
	rd.Period = p
	rd.Session = s
	rd.Tags = td
	rd.Region = r

	err := rd.GetAccountId()
	h.LogErrorExit(err)

	err = rd.BuildFilters()
	h.LogErrorExit(err)

	err = rd.CreateNamespaces()
	h.LogErrorExit(err)

	return nil
}

func (rd *RegionDescription) CreateNamespaces() error {
	namespaces := GetNamespaces()
	rd.Namespaces = make(map[string]*NamespaceDescription)
	for _, namespace := range namespaces {
		nd := NamespaceDescription{
			Namespace: aws.String(namespace),
			Parent:    rd,
		}
		rd.Namespaces[namespace] = &nd
	}

	return nil
}

func (rd *RegionDescription) GatherMetrics(cw *cloudwatch.CloudWatch) {
	log.Infof("Gathering metrics for region %s...", *rd.Region)
	rd.SetRequestTime()

	ndc := make(chan *NamespaceDescription)
	for _, namespace := range rd.Namespaces {
		// Initialize metric containers if they don't already exist
		for _, metric := range namespace.Metrics {
			for _, stat := range metric.Statistic {
				metric.initializeMetric(*stat)
			}
		}
		go namespace.GatherMetrics(cw, ndc)
	}
}

func (nd *NamespaceDescription) GatherMetrics(cw *cloudwatch.CloudWatch, ndc chan *NamespaceDescription) {
	for _, r := range nd.Resources {
		resource := r
		go func(rd *ResourceDescription, ndc chan *NamespaceDescription) {
			resource.Parent = nd
			result, err := resource.GetData(cw)
			h.LogError(err)
			err = resource.SaveData(result)
			h.LogError(err)
			ndc <- nd
		}(resource, ndc)
	}
}

func (m *MetricDescription) initializeMetric(stat string) {
	name := *m.metricName(stat)
	if _, ok := results[name]; ok == true {
		// metric is already initialized
		return
	}

	var promMetric prometheus.Collector
	if *m.Type == "counter" && (stat == "Sum" || stat == "SampleCount") {
		promMetric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: *m.Help,
			},
			[]string{"name", "id", "type", "region"},
		)
	} else {
		if *m.Type == "counter" {
			log.Debugf(
				"Cannot use metric type counter for stat %s. Metric %s will use a gauge instead",
				stat,
				name,
			)
		}
		promMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: *m.Help,
			},
			[]string{"name", "id", "type", "region"},
		)
	}
	results[name] = promMetric
	if err := prometheus.Register(promMetric); err != nil {
		log.Fatalf("Error registering metric %s: %s", name, err)
	}
}

func (rd *ResourceDescription) BuildDimensions(dd []*DimensionDescription) error {
	dl := []*cloudwatch.Dimension{}
	for _, dimension := range dd {
		c := &cloudwatch.Dimension{
			Name:  dimension.Name,
			Value: dimension.Value,
		}
		dl = append(dl, c)
	}
	rd.Dimensions = dl

	return nil
}

func (rd *ResourceDescription) BuildQuery() error {
	query := []*cloudwatch.MetricDataQuery{}
	for key, value := range rd.Parent.Metrics {
		dimensions := rd.Dimensions
		dimensions = append(dimensions, value.Dimensions...)
		for _, stat := range value.Statistic {
			cm := &cloudwatch.MetricDataQuery{
				Id: value.metricName(*stat),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String(key),
						Namespace:  rd.Parent.Namespace,
						Dimensions: dimensions,
					},
					Stat:   stat,
					Period: aws.Int64(int64(value.Period)),
				},
				// We hardcode the label so that we can rely on the ordering in
				// SaveData.
				Label:      aws.String(fmt.Sprintf("%s %s", key, *stat)),
				ReturnData: aws.Bool(true),
			}
			query = append(query, cm)
		}
	}
	rd.Query = query

	return nil
}

func (rd *ResourceDescription) SaveData(c *cloudwatch.GetMetricDataOutput) error {
	for _, data := range c.MetricDataResults {
		// Filter out old samples so they don't get double counted
		values := data.Values
		if len(values) <= 0 {
			continue
		}

		promLabels := prometheus.Labels{
			"name":   *rd.Name,
			"id":     *rd.ID,
			"type":   *rd.Type,
			"region": *rd.Parent.Parent.Region,
		}

		if counter, ok := results[*data.Id].(*prometheus.CounterVec); ok == true {
			key, err := counter.GetMetricWith(promLabels)
			if err != nil {
				return err
			}
			if lastTimestamp, ok := timestamps[key]; ok == true {
				values = h.NewValues(data.Values, data.Timestamps, *lastTimestamp)
			}
			if len(values) <= 0 {
				continue
			}
			// AWS returns the data in descending order
			timestamps[key] = data.Timestamps[0]
		}

		labels := strings.Split(*data.Label, " ")
		stat := labels[len(labels)-1]

		value := 0.0
		var err error = nil
		switch stat {
		case "Average":
			value, err = h.Average(values)
		case "Sum":
			value, err = h.Sum(values)
		case "Minimum":
			value, err = h.Min(values)
		case "Maximum":
			value, err = h.Max(values)
		case "SampleCount":
			value, err = h.Sum(values)
		default:
			err = fmt.Errorf("Unknown Statistic type: %s", stat)
		}
		if err != nil {
			return err
		}

		rd.Parent.Parent.Mutex.Lock()
		err = updateMetric(*data.Id, value, promLabels)
		rd.Parent.Parent.Mutex.Unlock()
		if err != nil {
			return err
		}
	}

	return nil
}

func updateMetric(name string, value float64, labels prometheus.Labels) error {
	if metric, ok := results[name]; ok == true {
		switch m := metric.(type) {
		case *prometheus.GaugeVec:
			m.With(labels).Set(value)
		case *prometheus.CounterVec:
			m.With(labels).Add(value)
		default:
			return fmt.Errorf("Could not resolve type of metric %s", name)
		}
	} else {
		return fmt.Errorf("Couldn't save metric %s", name)
	}
	return nil
}

func (rd *RegionDescription) TagsFound(tl interface{}) bool {
	tags := []*TagDescription{}

	// Not sure how to deal with code duplication here
	switch i := tl.(type) {
	case *elb.TagDescription:
		if len(i.Tags) < 1 {
			return false
		}
		for _, tag := range i.Tags {
			t := TagDescription{}
			awsutil.Copy(&t, tag)
			tags = append(tags, &t)
		}
	case *elbv2.TagDescription:
		if len(i.Tags) < 1 {
			return false
		}
		for _, tag := range i.Tags {
			t := TagDescription{}
			awsutil.Copy(&t, tag)
			tags = append(tags, &t)
		}
	case *rds.ListTagsForResourceOutput:
		if len(i.TagList) < 1 {
			return false
		}
		for _, tag := range i.TagList {
			t := TagDescription{}
			awsutil.Copy(&t, tag)
			tags = append(tags, &t)
		}
	case *elasticache.TagListMessage:
		if len(i.TagList) < 1 {
			return false
		}
		for _, tag := range i.TagList {
			t := TagDescription{}
			awsutil.Copy(&t, tag)
			tags = append(tags, &t)
		}
	case *s3.GetBucketTaggingOutput:
		if len(i.TagSet) < 1 {
			return false
		}
		for _, tag := range i.TagSet {
			t := TagDescription{}
			awsutil.Copy(&t, tag)
			tags = append(tags, &t)
		}
	default:
		return false
	}

	l1 := len(rd.Tags)
	l2 := len(tags)
	number_of_negative_matches := l1

	if l1 > l2 {
		return false
	}

	for _, left := range rd.Tags {
		for _, right := range tags {
			if *left.Key == *right.Key && *left.Value == *right.Value {
				number_of_negative_matches--
				break
			}
		}
		if number_of_negative_matches == 0 {
			return true
		}
	}
	return false
}

func (rd *ResourceDescription) GetData(cw *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	rd.BuildQuery()
	var startTime *time.Time

	if val, ok := rd.Parent.Parent.Config.Metrics.Data[*rd.Parent.Namespace]; ok {
		if val.Length > 0 {
			time := rd.Parent.Parent.Time.EndTime.Add(time.Minute * -time.Duration(val.Length))
			startTime = &time
		}
	} else {
		startTime = rd.Parent.Parent.Time.StartTime
	}
	input := cloudwatch.GetMetricDataInput{
		StartTime:         startTime,
		EndTime:           rd.Parent.Parent.Time.EndTime,
		MetricDataQueries: rd.Query,
	}
	result, err := cw.GetMetricData(&input)
	h.LogError(err)

	return result, err
}

func SameErrorType(l error, r *string) bool {
	if l != nil {
		if aerr, ok := l.(awserr.Error); ok {
			code := aerr.Code()
			if code == *r {
				return true
			}
		}
		return false
	}
	return false
}
