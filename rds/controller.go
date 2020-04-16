package rds

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

func createResourceDescription(nd *b.NamespaceDescription, dbi *rds.DBInstance) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("DBInstanceIdentifier"),
			Value: dbi.DBInstanceIdentifier,
		},
	}
	if err := rd.BuildDimensions(dd); err != nil {
		return nil, err
	}

	rd.ID = dbi.DBInstanceIdentifier
	rd.Name = dbi.DBInstanceIdentifier
	rd.Type = aws.String("rds")
	rd.Parent = nd

	return &rd, nil
}

// CreateResourceList fetches a list of all RDS databases in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating RDS resource list ...")
	session := rds.New(nd.Parent.Session)
	input := rds.DescribeDBInstancesInput{}
	result, err := session.DescribeDBInstances(&input)
	h.LogIfError(err)

	var w sync.WaitGroup
	w.Add(len(result.DBInstances))
	ch := make(chan *b.ResourceDescription, len(result.DBInstances))
	for _, dbi := range result.DBInstances {
		go func(dbi *rds.DBInstance, wg *sync.WaitGroup) {
			defer wg.Done()
			input := rds.ListTagsForResourceInput{
				ResourceName: dbi.DBInstanceArn,
			}
			tags, err := session.ListTagsForResource(&input)
			h.LogIfError(err)

			if nd.Parent.TagsFound(tags) {
				if r, err := createResourceDescription(nd, dbi); err == nil {
					ch <- r
				}
				h.LogIfError(err)
			}
		}(dbi, &w)
	}
	w.Wait()
	close(ch)

	resources := []*b.ResourceDescription{}
	for r := range ch {
		resources = append(resources, r)
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
