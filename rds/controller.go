package rds

import (
        b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
        h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, dbi *rds.DBInstance) error {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("DBInstanceIdentifier"),
			Value: dbi.DBInstanceIdentifier,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = dbi.DBInstanceIdentifier
	rd.Name = dbi.DBInstanceIdentifier
	rd.Type = aws.String("rds")
	rd.Parent = nd
	rd.BuildQuery()
	nd.Mutex.Lock()
	nd.Resources = append(nd.Resources, &rd)
	nd.Mutex.Unlock()

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating RDS resource list ...")
	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := rds.New(nd.Parent.Session)
	input := rds.DescribeDBInstancesInput{}
	result, err := session.DescribeDBInstances(&input)
	h.LogError(err)

	var w sync.WaitGroup
	w.Add(len(result.DBInstances))
	for _, dbi := range result.DBInstances {
		go func(dbi *rds.DBInstance, wg *sync.WaitGroup) {
			defer wg.Done()
			input := rds.ListTagsForResourceInput{
				ResourceName: dbi.DBInstanceArn,
			}
			tags, err := session.ListTagsForResource(&input)
			h.LogError(err)

			if nd.Parent.TagsFound(tags) {
				err := CreateResourceDescription(nd, dbi)
				h.LogError(err)
			}
		}(dbi, &w)
	}
	w.Wait()
	return nil
}
