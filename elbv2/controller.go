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
	lbID := strings.Split(*td.ResourceArn, "loadbalancer/")[1]
	lbTypeAndName := strings.Split(lbID, "/")
	lbName := lbTypeAndName[1]

	rd := b.ResourceDescription{}
	if lbTypeAndName[0] == "net" && *nd.Namespace == "AWS/NetworkELB" {
		rd.Type = aws.String("lb-network")
	} else if lbTypeAndName[0] == "app" && *nd.Namespace == "AWS/ApplicationELB" {
		rd.Type = aws.String("lb-application")
	} else {
		return nil
	}

	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("LoadBalancer"),
			Value: aws.String(lbID),
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = td.ResourceArn
	rd.Name = &lbName
	rd.Parent = nd
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

// CreateResourceList fetches a list of all ALB/NLB resources in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating ALB/NLB resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	session := elbv2.New(nd.Parent.Session)
	input := elbv2.DescribeLoadBalancersInput{}
	result, err := session.DescribeLoadBalancers(&input)
	h.LogError(err)

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
		h.LogError(err)
		tagDescriptions = append(tagDescriptions, tags.TagDescriptions...)
	}

	for _, td := range tagDescriptions {
		if nd.Parent.TagsFound(td) {
			err := CreateResourceDescription(nd, td)
			h.LogError(err)
		} else {
			continue
		}
	}
	return nil
}
