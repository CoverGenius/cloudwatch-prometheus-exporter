package network

import (
	"sync"

	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func createResourceDescription(nd *b.NamespaceDescription, ng *ec2.NatGateway) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("NatGatewayId"),
			Value: ng.NatGatewayId,
		},
	}
	if err := rd.BuildDimensions(dd); err != nil {
		return nil, err
	}

	rd.ID = ng.NatGatewayId
	rd.Name = ng.NatGatewayId
	rd.Type = aws.String("nat-gateway")
	rd.Parent = nd

	return &rd, nil
}

// CreateResourceList fetches a list of all NAT gateways in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating NatGateway resource list ...")
	session := ec2.New(nd.Parent.Session)
	input := ec2.DescribeNatGatewaysInput{
		Filter: nd.Parent.Filters,
	}
	result, err := session.DescribeNatGateways(&input)
	h.LogIfError(err)

	resources := []*b.ResourceDescription{}
	for _, ng := range result.NatGateways {
		if r, err := createResourceDescription(nd, ng); err == nil {
			resources = append(resources, r)
		}
		h.LogIfError(err)
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
