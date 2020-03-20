package elb

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, td *elb.TagDescription) error {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("LoadBalancerName"),
			Value: td.LoadBalancerName,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = td.LoadBalancerName
	rd.Name = td.LoadBalancerName
	rd.Type = aws.String("lb-classic")
	rd.Parent = nd
	rd.BuildQuery()
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating Classic LB resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := elb.New(nd.Parent.Session)
	input := elb.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogError(err)

	resource_list := []*string{}
	for _, lb := range result.LoadBalancerDescriptions {
		resource_list = append(resource_list, lb.LoadBalancerName)
	}
	dti := elb.DescribeTagsInput{
		LoadBalancerNames: resource_list,
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
