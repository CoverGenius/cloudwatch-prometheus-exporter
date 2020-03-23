package base

import (
	"fmt"
	"math"
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
	log "github.com/sirupsen/logrus"
)

var (
	nan     float64 = math.NaN()
	Results         = Metrics{
		Metric: make(map[string]map[string]*MetricDescription),
	}
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
	Statistic  *string
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
	Results    map[*NamespaceDescription]map[string][]*string
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
	rd.Results = make(map[*NamespaceDescription]map[string][]*string)
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
	rd.Results = map[*NamespaceDescription]map[string][]*string{}
	rd.SetRequestTime()

	ndc := make(chan *NamespaceDescription)
	for _, namespace := range rd.Namespaces {
		go namespace.GatherMetrics(cw, ndc)
	}

	for {
		select {
		case namespace := <-ndc:
			Results.Mutex.Lock()
			rd.Mutex.Lock()
			for metric := range rd.Results[namespace] {
				if _, ok := Results.Metric[*namespace.Namespace]; ok == false {
					Results.Metric[*namespace.Namespace] = make(map[string]*MetricDescription)
				}
				if _, ok := Results.Metric[*namespace.Namespace][metric]; ok == false {
					Results.Metric[*namespace.Namespace][metric] = &MetricDescription{
						Data:       make(map[string][]*string),
						Help:       namespace.Metrics[metric].Help,
						Type:       namespace.Metrics[metric].Type,
						OutputName: namespace.Metrics[metric].OutputName,
						Period:     namespace.Metrics[metric].Period,
						Statistic:  namespace.Metrics[metric].Statistic,
						Dimensions: namespace.Metrics[metric].Dimensions,
					}
				}
				if _, ok := Results.Metric[*namespace.Namespace][metric].Data[*rd.Region]; ok == false {
					Results.Metric[*namespace.Namespace][metric].Data[*rd.Region] = []*string{}
				}
				Results.Metric[*namespace.Namespace][metric].Data[*rd.Region] = rd.Results[namespace][metric]
			}
			rd.Mutex.Unlock()
			Results.Mutex.Unlock()
		}
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
		cm := &cloudwatch.MetricDataQuery{
			Id: rd.Parent.Metrics[key].OutputName,
			MetricStat: &cloudwatch.MetricStat{
				Metric: &cloudwatch.Metric{
					MetricName: aws.String(key),
					Namespace:  rd.Parent.Namespace,
					Dimensions: dimensions,
				},
				Stat:   value.Statistic,
				Period: aws.Int64(int64(value.Period)),
			},
			ReturnData: aws.Bool(true),
		}
		query = append(query, cm)
	}
	rd.Query = query

	return nil
}

func (rd *ResourceDescription) SaveData(c *cloudwatch.GetMetricDataOutput) error {
	for _, data := range c.MetricDataResults {
		size := float64(len(data.Values))
		value := &nan
		if size > 0 {
			average, err := h.CountAverage(data.Values, &size)
			h.LogError(err)
			value = average
		}
		result := fmt.Sprintf(
			"%s{name=\"%s\",id=\"%s\",type=\"%s\",region=\"%s\"} %.2f\n",
			*data.Id,
			*rd.Name,
			*rd.ID,
			*rd.Type,
			*rd.Parent.Parent.Region,
			*value,
		)

		rd.Parent.Parent.Mutex.Lock()
		labels := strings.Split(*data.Label, " ")
		metric := labels[len(labels)-1]
		if _, ok := rd.Parent.Parent.Results[rd.Parent]; ok == false {
			rd.Parent.Parent.Results[rd.Parent] = make(map[string][]*string)
		}
		rd.Parent.Parent.Results[rd.Parent][metric] = append(rd.Parent.Parent.Results[rd.Parent][metric], &result)
		rd.Parent.Parent.Mutex.Unlock()
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
