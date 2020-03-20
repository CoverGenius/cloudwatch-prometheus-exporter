package network

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, ng *ec2.NatGateway) error {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("NatGatewayId"),
			Value: ng.NatGatewayId,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)

	rd.ID = ng.NatGatewayId
	rd.Name = ng.NatGatewayId
	rd.Type = aws.String("nat-gateway")
	rd.Parent = nd
	rd.BuildQuery()
	nd.Resources = append(nd.Resources, &rd)

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating NatGateway resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := ec2.New(nd.Parent.Session)
	input := ec2.DescribeNatGatewaysInput{
		Filter: nd.Parent.Filters,
	}
	result, err := session.DescribeNatGateways(&input)
	h.LogError(err)

	for _, ng := range result.NatGateways {
		err := CreateResourceDescription(nd, ng)
		h.LogError(err)
	}
	return nil
}
