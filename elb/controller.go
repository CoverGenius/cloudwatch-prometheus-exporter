package elb

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
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
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

// CreateResourceList fetches a list of all Classic LB resources in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating Classic LB resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := elb.New(nd.Parent.Session)
	input := elb.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogError(err)

	resourceList := []*string{}
	for _, lb := range result.LoadBalancerDescriptions {
		resourceList = append(resourceList, lb.LoadBalancerName)
	}
	if len(resourceList) <= 0 {
		return nil
	}

	dti := elb.DescribeTagsInput{
		LoadBalancerNames: resourceList,
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
