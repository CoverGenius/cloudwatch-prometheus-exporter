package ec2

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createResourceDescription(nd *b.NamespaceDescription, instance *ec2.Instance) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("InstanceId"),
			Value: instance.InstanceId,
		},
	}
	if err := rd.BuildDimensions(dd); err != nil {
		return nil, err
	}

	tl := []string{}
	tags := make(map[string]*string)
	for _, t := range instance.Tags {
		tags[*t.Key] = t.Value
		ts := fmt.Sprintf("%s=%s", *t.Key, *t.Value)
		tl = append(tl, ts)
	}

	if len(tl) < 1 {
		rd.Tags = aws.String("")
	} else {
		rd.Tags = aws.String(strings.Join(tl, ","))
	}

	rd.ID = instance.InstanceId
	rd.Name = instance.InstanceId
	if name, ok := tags["Name"]; ok {
		rd.Name = name
	}
	rd.Type = aws.String("ec2")
	rd.Parent = nd

	return &rd, nil
}

// CreateResourceList fetches a list of all EC2 instances in the parent region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating EC2 resource list ...")

	session := ec2.New(nd.Parent.Session)
	input := ec2.DescribeInstancesInput{
		Filters: nd.Parent.Filters,
	}
	result, err := session.DescribeInstances(&input)
	h.LogIfError(err)

	resources := []*b.ResourceDescription{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if r, err := createResourceDescription(nd, instance); err == nil {
				resources = append(resources, r)
			}
			h.LogIfError(err)
		}
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
