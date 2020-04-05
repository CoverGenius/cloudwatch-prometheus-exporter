package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/ec2"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elasticache"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elb"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/elbv2"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/network"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/rds"
	"github.com/CoverGenius/cloudwatch-prometheus-exporter/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	log "github.com/sirupsen/logrus"
)

var (
	rdd    []*base.RegionDescription
	config string
)

func init() {
	flag.StringVar(&config, "config", "config.yaml", "Path to config file")
}

func run(nd map[string]*base.NamespaceDescription, cw *cloudwatch.CloudWatch, rd *base.RegionDescription, pi uint8) {
	var delay uint8 = 0
	for {
		select {
		case <-time.After(time.Duration(delay) * time.Minute):
			var wg sync.WaitGroup
			wg.Add(8)
			log.Debug("Creating list of resources ...")
			go elasticache.CreateResourceList(nd["AWS/ElastiCache"], &wg)
			go rds.CreateResourceList(nd["AWS/RDS"], &wg)
			go ec2.CreateResourceList(nd["AWS/EC2"], &wg)
			go network.CreateResourceList(nd["AWS/NATGateway"], &wg)
			go elb.CreateResourceList(nd["AWS/ELB"], &wg)
			go elbv2.CreateResourceList(nd["AWS/ApplicationELB"], &wg)
			go elbv2.CreateResourceList(nd["AWS/NetworkELB"], &wg)
			go s3.CreateResourceList(nd["AWS/S3"], &wg)
			wg.Wait()
			log.Debug("Gathering metrics ...")
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

	if c.Period == 0 {
		c.Period = 5
	}

	if c.APIKey == "" || c.APISecret == "" {
		log.Fatal("Either api_key or api_secret is empty!")
	}

	if len(c.Regions) < 1 {
		log.Fatal("No regions specified. Please set at least one!")
	}

	if c.PollInterval == 0 {
		c.PollInterval = 5
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(h.GetLogLevel(c.LogLevel))

	os.Setenv("AWS_ACCESS_KEY_ID", c.APIKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", c.APISecret)

	return &c
}

func main() {
	flag.Parse()
	c := processConfig(&config)

	for _, region := range c.Regions {
		r := region
		session := session.Must(session.NewSession(&aws.Config{Region: r}))
		cw := cloudwatch.New(session)
		rd := base.RegionDescription{
			Config: c,
		}
		rdd = append(rdd, &rd)
		rd.Init(session, c.Tags, r, &c.Period)

		go run(rd.Namespaces, cw, &rd, c.PollInterval)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(c.Listen, nil))
}
