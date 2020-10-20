package elbv2

import (
	"errors"

	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	log "github.com/sirupsen/logrus"

	"strings"
	"sync"
)

func createResourceDescription(nd *b.NamespaceDescription, tags *string, td *elbv2.TagDescription) (*b.ResourceDescription, error) {
	lbID := strings.Split(*td.ResourceArn, "loadbalancer/")[1]
	lbTypeAndName := strings.Split(lbID, "/")
	lbName := lbTypeAndName[1]

	rd := b.ResourceDescription{}
	switch {
	case lbTypeAndName[0] == "net" && *nd.Namespace == "AWS/NetworkELB":
		rd.Type = aws.String("lb-network")
	case lbTypeAndName[0] == "app" && *nd.Namespace == "AWS/ApplicationELB":
		rd.Type = aws.String("lb-application")
	default:
		return nil, errors.New("invalid lb type")
	}

	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("LoadBalancer"),
			Value: aws.String(lbID),
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogIfError(err)
	rd.ID = td.ResourceArn
	rd.Name = &lbName
	rd.Parent = nd
	rd.Tags = tags

	return &rd, nil
}

// CreateResourceList fetches a list of all ALB/NLB resources in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating ALB/NLB resource list ...")
	session := elbv2.New(nd.Parent.Session)
	input := elbv2.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogIfError(err)

	resourceList := []*string{}
	for _, lb := range result.LoadBalancers {
		resourceList = append(resourceList, lb.LoadBalancerArn)
	}

	// The AWS ELBV2 API has a limit of 20 resources which can be described in one request
	chunkSize := 20
	tagDescriptions := []*elbv2.TagDescription{}
	for i := 0; i < len(resourceList); i += chunkSize {
		end := i + chunkSize
		if end > len(resourceList) {
			end = len(resourceList)
		}

		dti := elbv2.DescribeTagsInput{
			ResourceArns: resourceList[i:end],
		}
		tags, err := session.DescribeTags(&dti)
		h.LogIfError(err)
		tagDescriptions = append(tagDescriptions, tags.TagDescriptions...)
	}

	resources := []*b.ResourceDescription{}
	for _, td := range tagDescriptions {
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
