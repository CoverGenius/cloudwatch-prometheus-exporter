package elasticache

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"strings"
	"sync"
)

func CreateResourceDescription(nd *b.NamespaceDescription, cc *elasticache.CacheCluster) error {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("CacheClusterId"),
			Value: cc.CacheClusterId,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogError(err)
	rd.ID = cc.CacheClusterId
	rd.Name = cc.CacheClusterId
	rd.Type = aws.String("elasticache")
	rd.Parent = nd
	rd.BuildQuery()
	nd.Mutex.Lock()
	nd.Resources = append(nd.Resources, &rd)
	nd.Mutex.Unlock()

	return nil
}

func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debug("Creating Elasticache resource list ...")

	nd.Resources = []*b.ResourceDescription{}
	nd.Metrics = GetMetrics()
	session := elasticache.New(nd.Parent.Session)
	input := elasticache.DescribeCacheClustersInput{}
	result, err := session.DescribeCacheClusters(&input)
	h.LogError(err)
	service := "elasticache"

	var w sync.WaitGroup
	w.Add(len(result.CacheClusters))
	for _, cc := range result.CacheClusters {
		go func(cc *elasticache.CacheCluster, wg *sync.WaitGroup) {
			defer wg.Done()
			resource := strings.Join([]string{"cluster", *cc.CacheClusterId}, ":")
			arn, err := nd.Parent.BuildArn(&service, &resource)
			h.LogError(err)
			input := elasticache.ListTagsForResourceInput{
				ResourceName: aws.String(arn),
			}
			tags, err := session.ListTagsForResource(&input)
			h.LogError(err)

			if nd.Parent.TagsFound(tags) {
				err := CreateResourceDescription(nd, cc)
				h.LogError(err)
			}
		}(cc, &w)
	}
	w.Wait()
	return nil
}
