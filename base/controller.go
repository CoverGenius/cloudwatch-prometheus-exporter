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
	results = make(map[string]prometheus.Collector)
)

// TimeDescription represents an interval with a specific start and finish time
type TimeDescription struct {
	StartTime *time.Time
	EndTime   *time.Time
}

// TagDescription represents an AWS tag key value pair
type TagDescription struct {
	Key   *string `yaml:"name"`
	Value *string `yaml:"value"`
}

// DimensionDescription represents a Cloudwatch dimension key value pair
type DimensionDescription struct {
	Name  *string
	Value *string
}

// MetricDescription describes a single Cloudwatch metric with one or more
// statistics to be monitored for relevant resources
type MetricDescription struct {
	Help       *string
	OutputName *string
	Dimensions []*cloudwatch.Dimension
	Period     int
	Statistic  []*string
	Data       map[string][]*string
}

// RegionDescription describes an AWS region which will be monitored via cloudwatch
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

// NamespaceDescription describes an AWS namespace to be monitored via cloudwatch
// e.g. EC2 or S3
type NamespaceDescription struct {
	Namespace *string
	Resources []*ResourceDescription
	Parent    *RegionDescription
	Mutex     sync.RWMutex
	Metrics   map[string]*MetricDescription
}

// ResourceDescription describes a single AWS resource which will be monitored via
// one or more cloudwatch metrics.
type ResourceDescription struct {
	Name       *string
	ID         *string
	Dimensions []*cloudwatch.Dimension
	Type       *string
	Parent     *NamespaceDescription
	Mutex      sync.RWMutex
	Query      []*cloudwatch.MetricDataQuery
	timestamps map[prometheus.Collector]*time.Time
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

func (rd *RegionDescription) setRequestTime() error {
	log.Debug("Setting request time ...")
	td := TimeDescription{}
	t := time.Now().Round(time.Minute * 5)
	start := t.Add(time.Minute * -time.Duration(*rd.Period))
	td.StartTime = &start
	td.EndTime = &t
	rd.Time = &td
	return nil
}

// BuildARN returns the AWS ARN of a resource in a region given the input service and resource
func (rd *RegionDescription) BuildARN(s *string, r *string) (string, error) {
	a := arn.ARN{
		Service:   *s,
		Region:    *rd.Region,
		AccountID: *rd.AccountID,
		Resource:  *r,
		Partition: "aws",
	}
	return a.String(), nil
}

func (rd *RegionDescription) buildFilters() error {
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

func (rd *RegionDescription) saveAccountID() error {
	session := iam.New(rd.Session)
	input := iam.GetUserInput{}
	user, err := session.GetUser(&input)
	h.LogError(err)
	a, err := arn.Parse(*user.User.Arn)
	h.LogError(err)
	rd.AccountID = &a.AccountID

	return nil
}

// Init initializes a region and its nested namspaces in preparation for collection
// cloudwatchc metrics for that region.
func (rd *RegionDescription) Init(s *session.Session, td []*TagDescription, r *string, p *uint8) error {
	log.Infof("Initializing region %s ...", *r)
	rd.Period = p
	rd.Session = s
	rd.Tags = td
	rd.Region = r

	err := rd.saveAccountID()
	h.LogErrorExit(err)

	err = rd.buildFilters()
	h.LogErrorExit(err)

	err = rd.CreateNamespaceDescriptions()
	h.LogErrorExit(err)

	return nil
}

// CreateNamespaceDescriptions populates the list of NamespaceDescriptions for an AWS region
func (rd *RegionDescription) CreateNamespaceDescriptions() error {
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

// GatherMetrics queries the Cloudwatch API for metrics related to the resources in this region
func (rd *RegionDescription) GatherMetrics(cw *cloudwatch.CloudWatch) {
	log.Infof("Gathering metrics for region %s...", *rd.Region)
	rd.setRequestTime()

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

// GatherMetrics queries the Cloudwatch API for metrics related to this AWS namespace in the parent region
func (nd *NamespaceDescription) GatherMetrics(cw *cloudwatch.CloudWatch, ndc chan *NamespaceDescription) {
	for _, r := range nd.Resources {
		resource := r
		go func(rd *ResourceDescription, ndc chan *NamespaceDescription) {
			resource.Parent = nd
			result, err := resource.getData(cw)
			h.LogError(err)
			resource.saveData(result)
			ndc <- nd
		}(resource, ndc)
	}
}

func (md *MetricDescription) initializeMetric(stat string) {
	name := *md.metricName(stat)
	if _, ok := results[name]; ok == true {
		// metric is already initialized
		return
	}

	var promMetric prometheus.Collector
	if stat == "Sum" || stat == "SampleCount" {
		promMetric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: *md.Help,
			},
			[]string{"name", "id", "type", "region"},
		)
	} else {
		promMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: *md.Help,
			},
			[]string{"name", "id", "type", "region"},
		)
	}
	results[name] = promMetric
	if err := prometheus.Register(promMetric); err != nil {
		log.Fatalf("Error registering metric %s: %s", name, err)
	}
}

// BuildDimensions coverts a slice of DimensionDescription to a slice of cloudwatchDimension and associates it with the resource
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

// BuildQuery constructs and saves the cloudwatch query for all the metrics associated with the resource
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
				// saveData.
				Label:      aws.String((&awsLabels{key, *stat}).String()),
				ReturnData: aws.Bool(true),
			}
			query = append(query, cm)
		}
	}
	rd.Query = query

	return nil
}

type awsLabels struct {
	metric    string
	statistic string
}

func (l *awsLabels) String() string {
	return fmt.Sprintf("%s %s", l.metric, l.statistic)
}

func awsLabelsFromString(s string) (*awsLabels, error) {
	stringLabels := strings.Split(s, " ")
	if len(stringLabels) < 2 {
		return nil, fmt.Errorf("Expected at least two labels, got %s", s)
	}
	labels := awsLabels{
		metric:    stringLabels[len(stringLabels)-2],
		statistic: stringLabels[len(stringLabels)-1],
	}
	return &labels, nil
}

func (rd *ResourceDescription) saveData(c *cloudwatch.GetMetricDataOutput) {
	for _, data := range c.MetricDataResults {
		if len(data.Values) <= 0 {
			continue
		}

		values, err := rd.filterValues(data)
		if len(values) <= 0 || err != nil {
			h.LogError(err)
			continue
		}

		labels, err := awsLabelsFromString(*data.Label)
		if err != nil {
			h.LogError(err)
			continue
		}

		value := 0.0
		switch labels.statistic {
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
			err = fmt.Errorf("Unknown Statistic type: %s", labels.statistic)
		}
		if err != nil {
			h.LogError(err)
			continue
		}

		err = rd.updateMetric(*data.Id, value)
		if err != nil {
			h.LogError(err)
			continue
		}
	}
}

func (rd *ResourceDescription) filterValues(data *cloudwatch.MetricDataResult) ([]*float64, error) {
	// In the case of a counter we need to remove any datapoints which have
	// already been added to the counter, otherwise if the poll intervals
	// overlap we will double count some data.
	values := data.Values
	if counter, ok := results[*data.Id].(*prometheus.CounterVec); ok == true {
		counter, err := counter.GetMetricWith(prometheus.Labels{
			"name":   *rd.Name,
			"id":     *rd.ID,
			"type":   *rd.Type,
			"region": *rd.Parent.Parent.Region,
		})
		if err != nil {
			return nil, err
		}
		rd.Mutex.Lock()
		defer rd.Mutex.Unlock()
		if rd.timestamps == nil {
			rd.timestamps = make(map[prometheus.Collector]*time.Time)
		}
		if lastTimestamp, ok := rd.timestamps[counter]; ok == true {
			values = h.NewValues(data.Values, data.Timestamps, *lastTimestamp)
		}
		if len(values) > 0 {
			// AWS returns the data in descending order
			rd.timestamps[counter] = data.Timestamps[0]
		}
	}
	return values, nil
}

func (rd *ResourceDescription) updateMetric(name string, value float64) error {
	labels := prometheus.Labels{
		"name":   *rd.Name,
		"id":     *rd.ID,
		"type":   *rd.Type,
		"region": *rd.Parent.Parent.Region,
	}
	rd.Parent.Parent.Mutex.Lock()
	defer rd.Parent.Parent.Mutex.Unlock()
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
	numberOfNegativeMatches := l1

	if l1 > l2 {
		return false
	}

	for _, left := range rd.Tags {
		for _, right := range tags {
			if *left.Key == *right.Key && *left.Value == *right.Value {
				numberOfNegativeMatches--
				break
			}
		}
		if numberOfNegativeMatches == 0 {
			return true
		}
	}
	return false
}

func (rd *ResourceDescription) getData(cw *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	rd.BuildQuery()

	startTime := rd.Parent.Parent.Time.StartTime
	if val, ok := rd.Parent.Parent.Config.Metrics.Data[*rd.Parent.Namespace]; ok {
		// Some resources don't have any data for a while (e.g. S3), in these cases the Period parameter
		// can be used to override the window used when querying the Cloudwatch API.
		if val.Period > 0 {
			time := rd.Parent.Parent.Time.EndTime.Add(-time.Duration(val.Period) * time.Minute)
			startTime = &time
		}
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

// IsSameErrorType returns true if the input error l is an AWS error with the same code as the input code r
func IsSameErrorType(l error, r *string) bool {
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
