package ec2

import (
        b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
        h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, instance *ec2.Instance) error {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("InstanceId"),
			Value: instance.InstanceId,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = instance.InstanceId
	rd.Name = instance.InstanceId
	rd.Type = aws.String("ec2")
	rd.Parent = nd
	rd.BuildQuery()
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating EC2 resource list ...")

	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := ec2.New(nd.Parent.Session)
	input := ec2.DescribeInstancesInput{
		Filters: nd.Parent.Filters,
	}
	result, err := session.DescribeInstances(&input)
	h.LogError(err)

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			err := CreateResourceDescription(nd, instance)
			h.LogError(err)
		}
	}
	return nil
}
