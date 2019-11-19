package elbv2

import (
        b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
        h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	log "github.com/sirupsen/logrus"

	"strings"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, td *elbv2.TagDescription) error {
	lb_id := strings.Split(*td.ResourceArn, "loadbalancer/")[1]
	lb_type_and_name := strings.Split(lb_id, "/")
	lb_name := lb_type_and_name[1]

	rd := b.ResourceDescription{}
	if lb_type_and_name[0] == "net" && *nd.Namespace == "AWS/NetworkELB" {
		rd.Type = aws.String("lb-network")
	} else if lb_type_and_name[0] == "app" && *nd.Namespace == "AWS/ApplicationELB" {
		rd.Type = aws.String("lb-application")
	} else {
		return nil
	}

	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("LoadBalancer"),
			Value: aws.String(lb_id),
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = td.ResourceArn
	rd.Name = &lb_name
	rd.Parent = nd
	rd.BuildQuery()
	nd.Metrics = GetMetrics(rd.Type)
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating ALB/NLB resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	session := elbv2.New(nd.Parent.Session)
	input := elbv2.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogError(err)

	resource_list := []*string{}
	for _, lb := range result.LoadBalancers {
		resource_list = append(resource_list, lb.LoadBalancerArn)
	}
	dti := elbv2.DescribeTagsInput{
		ResourceArns: resource_list,
	}
	tags, err := session.DescribeTags(&dti)
	h.LogError(err)

	for _, td := range tags.TagDescriptions {
		if nd.Parent.TagsFound(td) {
			err := CreateResourceDescription(nd, td)
			h.LogError(err)
		} else {
			continue
		}
	}
	return nil
}
