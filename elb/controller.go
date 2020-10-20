package elb

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
)

func createResourceDescription(nd *b.NamespaceDescription, tags *string, td *elb.TagDescription) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("LoadBalancerName"),
			Value: td.LoadBalancerName,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogIfError(err)
	rd.ID = td.LoadBalancerName
	rd.Name = td.LoadBalancerName
	rd.Type = aws.String("lb-classic")
	rd.Parent = nd
	rd.Tags = tags

	return &rd, nil
}

// CreateResourceList fetches a list of all Classic LB resources in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating Classic LB resource list ...")
	session := elb.New(nd.Parent.Session)
	input := elb.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogIfError(err)

	resourceList := []*string{}
	for _, lb := range result.LoadBalancerDescriptions {
		resourceList = append(resourceList, lb.LoadBalancerName)
	}
	if len(resourceList) <= 0 {
		return
	}

	dti := elb.DescribeTagsInput{
		LoadBalancerNames: resourceList,
	}
	tags, err := session.DescribeTags(&dti)
	h.LogIfError(err)

	resources := []*b.ResourceDescription{}
	for _, td := range tags.TagDescriptions {
		tl, found := nd.Parent.TagsFound(td)
		ts := b.TagsToString(tl)
		if found {
			if r, err := createResourceDescription(nd, ts, td); err == nil {
				resources = append(resources, r)
			}
			h.LogIfError(err)
		} else {
			continue
		}
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
