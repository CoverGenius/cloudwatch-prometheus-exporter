package base

import (
	"fmt"
	"regexp"
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

var alphaRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// TagDescription represents an AWS tag key value pair
type TagDescription struct {
	Key   *string `yaml:"name"`
	Value *string `yaml:"value"`
}

// DimensionDescription represents a Cloudwatch dimension key value pair
type DimensionDescription struct {
	Name  *string `yaml:"name"`
	Value *string `yaml:"value"`
}

// MetricDescription describes a single Cloudwatch metric with one or more
// statistics to be monitored for relevant resources
type MetricDescription struct {
	AWSMetric     string
	Namespace     string
	Help          *string
	OutputName    *string
	Dimensions    []*cloudwatch.Dimension
	PeriodSeconds int64
	RangeSeconds  int64
	Statistic     []*string

	timestamps map[awsLabels]*time.Time
	mutex      sync.RWMutex
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
	Mutex      sync.RWMutex
}

// NamespaceDescription describes an AWS namespace to be monitored via cloudwatch
// e.g. EC2 or S3
type NamespaceDescription struct {
	Namespace *string
	Resources []*ResourceDescription
	Parent    *RegionDescription
	Mutex     sync.RWMutex
	Metrics   []*MetricDescription
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

func (rd *RegionDescription) saveFilters() {
	filters := []*ec2.Filter{}
	for _, tag := range rd.Tags {
		f := &ec2.Filter{
			Name:   aws.String(strings.Join([]string{"tag", *tag.Key}, ":")),
			Values: []*string{tag.Value},
		}
		filters = append(filters, f)
	}
	rd.Filters = filters
}

func (rd *RegionDescription) saveAccountID() error {
	session := iam.New(rd.Session)
	input := iam.GetUserInput{}
	user, err := session.GetUser(&input)
	if err != nil {
		return err
	}

	a, err := arn.Parse(*user.User.Arn)
	if err != nil {
		return err
	}
	rd.AccountID = &a.AccountID

	return nil
}

// Init initializes a region and its nested namspaces in preparation for collection
// cloudwatchc metrics for that region.
func (rd *RegionDescription) Init(s *session.Session, td []*TagDescription, metrics map[string][]*MetricDescription) error {
	log.Infof("Initializing region %s ...", *rd.Region)
	rd.Session = s
	rd.Tags = td

	err := rd.saveAccountID()
	if err != nil {
		return fmt.Errorf("error saving account id: %s", err)
	}

	rd.saveFilters()

	err = rd.CreateNamespaceDescriptions(metrics)
	if err != nil {
		return fmt.Errorf("error creating namespaces: %s", err)
	}

	return nil
}

// CreateNamespaceDescriptions populates the list of NamespaceDescriptions for an AWS region
func (rd *RegionDescription) CreateNamespaceDescriptions(metrics map[string][]*MetricDescription) error {
	namespaces := GetNamespaces()
	rd.Namespaces = make(map[string]*NamespaceDescription)
	for _, namespace := range namespaces {
		nd := NamespaceDescription{
			Namespace: aws.String(namespace),
			Parent:    rd,
		}
		if mds, ok := metrics[namespace]; ok {
			nd.Metrics = mds
		}
		rd.Namespaces[namespace] = &nd
	}

	return nil
}

// GatherMetrics queries the Cloudwatch API for metrics related to the resources in this region
func (rd *RegionDescription) GatherMetrics(cw *cloudwatch.CloudWatch) {
	log.Infof("Gathering metrics for region %s...", *rd.Region)

	for _, namespace := range rd.Namespaces {
		go namespace.GatherMetrics(cw)
	}
}

// GatherMetrics queries the Cloudwatch API for metrics related to this AWS namespace in the parent region
func (nd *NamespaceDescription) GatherMetrics(cw *cloudwatch.CloudWatch) {
	for _, md := range nd.Metrics {
		go func(md *MetricDescription) {
			nd.Mutex.RLock()
			result, err := md.getData(cw, nd.Resources)
			nd.Mutex.RUnlock()
			h.LogIfError(err)
			md.saveData(result, *nd.Parent.Region)
		}(md)
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

func (rd *ResourceDescription) queryID(stat string) *string {
	// Cloudwatch calls need a snake-case-unique-id
	id := strings.ToLower(*rd.ID + "-" + stat)
	id = alphaRegex.ReplaceAllString(id, "_")
	return aws.String(id)
}

// BuildQuery constructs and saves the cloudwatch query for all the metrics associated with the resource
func (md *MetricDescription) BuildQuery(rds []*ResourceDescription) ([]*cloudwatch.MetricDataQuery, error) {
	query := []*cloudwatch.MetricDataQuery{}
	for _, rd := range rds {
		dimensions := rd.Dimensions
		dimensions = append(dimensions, md.Dimensions...)
		for _, stat := range md.Statistic {
			cm := &cloudwatch.MetricDataQuery{
				Id: rd.queryID(*stat),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: &md.AWSMetric,
						Namespace:  rd.Parent.Namespace,
						Dimensions: dimensions,
					},
					Stat:   stat,
					Period: aws.Int64(md.PeriodSeconds),
				},
				// We hardcode the label so that we can rely on the ordering in
				// saveData.
				Label:      aws.String((&awsLabels{*stat, *rd.Name, *rd.ID, *rd.Type, *rd.Parent.Parent.Region}).String()),
				ReturnData: aws.Bool(true),
			}
			query = append(query, cm)
		}
	}
	return query, nil
}

type awsLabels struct {
	statistic string
	name      string
	id        string
	rType     string
	region    string
}

func (l *awsLabels) String() string {
	return fmt.Sprintf("%s %s %s %s %s", l.statistic, l.name, l.id, l.rType, l.region)
}

func awsLabelsFromString(s string) (*awsLabels, error) {
	stringLabels := strings.Split(s, " ")
	if len(stringLabels) < 5 {
		return nil, fmt.Errorf("expected at least five labels, got %s", s)
	}
	labels := awsLabels{
		statistic: stringLabels[len(stringLabels)-5],
		name:      stringLabels[len(stringLabels)-4],
		id:        stringLabels[len(stringLabels)-3],
		rType:     stringLabels[len(stringLabels)-2],
		region:    stringLabels[len(stringLabels)-1],
	}
	return &labels, nil
}

func (md *MetricDescription) saveData(c *cloudwatch.GetMetricDataOutput, region string) {
	newData := map[string][]*promMetric{}
	for _, stat := range md.Statistic {
		// pre-allocate in case the last resource for a stat goes away
		newData[*stat] = []*promMetric{}
	}
	for _, data := range c.MetricDataResults {
		if len(data.Values) <= 0 {
			continue
		}

		labels, err := awsLabelsFromString(*data.Label)
		if err != nil {
			h.LogIfError(err)
			continue
		}

		values := md.filterValues(data, labels)
		if len(values) <= 0 {
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
			err = fmt.Errorf("unknown statistic type: %s", labels.statistic)
		}
		if err != nil {
			h.LogIfError(err)
			continue
		}

		newData[labels.statistic] = append(newData[labels.statistic], &promMetric{value, []string{labels.name, labels.id, labels.rType, labels.region}})
	}
	for stat, data := range newData {
		name := *md.metricName(stat)
		opts := prometheus.Opts{
			Name: name,
			Help: *md.Help,
		}
		labels := []string{"name", "id", "type", "region"}

		exporter.mutex.Lock()
		if _, ok := exporter.data[name+region]; !ok {
			if stat == "Sum" {
				exporter.data[name+region] = NewBatchCounterVec(opts, labels)
			} else {
				exporter.data[name+region] = NewBatchGaugeVec(opts, labels)
			}
		}
		exporter.mutex.Unlock()

		exporter.mutex.RLock()
		exporter.data[name+region].BatchUpdate(data)
		exporter.mutex.RUnlock()
	}
}

func (md *MetricDescription) filterValues(data *cloudwatch.MetricDataResult, labels *awsLabels) []*float64 {
	// In the case of a counter we need to remove any datapoints which have
	// already been added to the counter, otherwise if the poll intervals
	// overlap we will double count some data.
	values := data.Values
	if labels.statistic == "Sum" {
		md.mutex.Lock()
		defer md.mutex.Unlock()
		if md.timestamps == nil {
			md.timestamps = make(map[awsLabels]*time.Time)
		}
		if lastTimestamp, ok := md.timestamps[*labels]; ok {
			values = h.NewValues(data.Values, data.Timestamps, *lastTimestamp)
		}
		if len(values) > 0 {
			// AWS returns the data in descending order
			md.timestamps[*labels] = data.Timestamps[0]
		}
	}
	return values
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

func (md *MetricDescription) getData(cw *cloudwatch.CloudWatch, rds []*ResourceDescription) (*cloudwatch.GetMetricDataOutput, error) {
	query, err := md.BuildQuery(rds)
	if len(query) == 0 {
		return &cloudwatch.GetMetricDataOutput{}, nil
	}
	h.LogIfError(err)

	end := time.Now().Round(5 * time.Minute)
	start := end.Add(-time.Duration(md.RangeSeconds) * time.Second)

	input := cloudwatch.GetMetricDataInput{
		StartTime:         &start,
		EndTime:           &end,
		MetricDataQueries: query,
	}
	result, err := cw.GetMetricData(&input)
	h.LogIfError(err)

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
