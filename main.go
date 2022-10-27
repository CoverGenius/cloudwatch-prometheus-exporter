package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/CoverGenius/cloudwatch-prometheus-exporter/backup"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/ec2"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elasticache"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elb"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elbv2"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/network"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/s3"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/sqs"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/vpc"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/rds"
	log "github.com/sirupsen/logrus"
)

var (
	rdd    []*base.RegionDescription
	config string
)

func init() {
	flag.StringVar(&config, "config", "config.yaml", "Path to config file")
}

func run(nd map[string]*base.NamespaceDescription, cw *cloudwatch.CloudWatch, rd *base.RegionDescription, pi int64) {
	var delay int64 = 0
	for {
		select {
		case <-time.After(time.Duration(delay) * time.Second):
			var wg sync.WaitGroup
			wg.Add(11)
			log.Debug("Creating list of resources ...")
			go elasticache.CreateResourceList(nd["AWS/ElastiCache"], &wg)
			go rds.CreateResourceList(nd["AWS/RDS"], &wg)
			go ec2.CreateResourceList(nd["AWS/EC2"], &wg)
			go network.CreateResourceList(nd["AWS/NATGateway"], &wg)
			go elb.CreateResourceList(nd["AWS/ELB"], &wg)
			go elbv2.CreateResourceList(nd["AWS/ApplicationELB"], &wg)
			go elbv2.CreateResourceList(nd["AWS/NetworkELB"], &wg)
			go s3.CreateResourceList(nd["AWS/S3"], &wg)
			go sqs.CreateResourceList(nd["AWS/SQS"], &wg)
			go vpc.CreateResourceList(nd["AWS/VPC"], &wg)
			go backup.CreateResourceList(nd["AWS/Backup"], &wg)
			wg.Wait()
			delay = pi
			go rd.GatherMetrics(cw)
		}
	}
}

func processConfig(p *string) *base.Config {
	c := base.Config{}
	h.YAMLDecode(p, &c)

	if c.Listen == "" {
		c.Listen = "127.0.0.1:8080"
	}

	if c.PeriodSeconds == 0 {
		c.PeriodSeconds = 60
	}

	if c.RangeSeconds == 0 {
		c.RangeSeconds = 300
	}

	if c.APIKey == "" || c.APISecret == "" {
		log.Fatal("Either api_key or api_secret is empty!")
	}

	if len(c.Regions) < 1 {
		log.Fatal("No regions specified. Please set at least one!")
	}

	if c.PollInterval == 0 {
		c.PollInterval = 300
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(h.GetLogLevel(c.LogLevel))

	os.Setenv("AWS_ACCESS_KEY_ID", c.APIKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", c.APISecret)

	return &c
}

func main() {
	flag.Parse()
	// TODO allow hot reload of config
	c := processConfig(&config)
	defaults := map[string]map[string]*base.MetricDescription{
		"AWS/RDS":            rds.Metrics,
		"AWS/ElastiCache":    elasticache.Metrics,
		"AWS/EC2":            ec2.Metrics,
		"AWS/NATGateway":     network.Metrics,
		"AWS/ELB":            elb.Metrics,
		"AWS/ApplicationELB": elbv2.ALBMetrics,
		"AWS/NetworkELB":     elbv2.NLBMetrics,
		"AWS/S3":             s3.Metrics,
		"AWS/SQS":            sqs.Metrics,
		"AWS/VPC":            vpc.Metrics,
		"AWS/Backup":         backup.Metrics,
	}
	mds := c.ConstructMetrics(defaults)

	for _, r := range c.Regions {
		awsSession := session.Must(session.NewSession(&aws.Config{Region: r}))
		cw := cloudwatch.New(awsSession)
		rd := base.RegionDescription{Region: r}
		rdd = append(rdd, &rd)
		if err := rd.Init(awsSession, c.Tags, mds); err != nil {
			log.Fatalf("error initializing region: %s", err)
		}

		go run(rd.Namespaces, cw, &rd, c.PollInterval)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(c.Listen, nil))
}
