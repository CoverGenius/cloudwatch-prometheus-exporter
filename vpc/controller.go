package vpc

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createResourceDescription(nd *b.NamespaceDescription, subnet *ec2.Subnet) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("SubnetId"),
			Value: subnet.SubnetId,
		},
	}
	if err := rd.BuildDimensions(dd); err != nil {
		return nil, err
	}

	tl := []string{}
	tags := make(map[string]*string)
	for _, t := range subnet.Tags {
		*t.Value = strings.ReplaceAll(*t.Value, " ", "_")
		tags[*t.Key] = t.Value
		ts := fmt.Sprintf("%s=%s", *t.Key, *t.Value)
		tl = append(tl, ts)
	}

	sort.Strings(tl)

	if len(tl) < 1 {
		rd.Tags = aws.String("")
	} else {
		rd.Tags = aws.String(strings.Join(tl, ","))
	}

	rd.ID = subnet.SubnetId
	rd.Name = subnet.SubnetId
	rd.Type = aws.String("vpc")
	rd.Parent = nd

	return &rd, nil
}

// CreateResourceList fetches a list of all EC2 instances in the parent region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating subnet resource list ...")

	session := ec2.New(nd.Parent.Session)
	input := ec2.DescribeSubnetsInput{
		Filters: nd.Parent.Filters,
	}
	result, err := session.DescribeSubnets(&input)
	h.LogIfError(err)

	resources := []*b.ResourceDescription{}
	for _, subnet := range result.Subnets {
		if r, err := createResourceDescription(nd, subnet); err == nil {
			resources = append(resources, r)
		}
		h.LogIfError(err)
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
