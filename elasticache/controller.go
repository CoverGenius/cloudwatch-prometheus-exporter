package elasticache

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

func createResourceDescription(nd *b.NamespaceDescription, tags *string, cc *elasticache.CacheCluster) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("CacheClusterId"),
			Value: cc.CacheClusterId,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogIfError(err)
	rd.ID = cc.CacheClusterId
	rd.Name = cc.CacheClusterId
	rd.Type = aws.String("elasticache")
	rd.Parent = nd
	rd.Tags = tags

	return &rd, nil
}

// CreateResourceList fetches a list of all Elasticache clusters in the parent region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating Elasticache resource list ...")

	session := elasticache.New(nd.Parent.Session)
	input := elasticache.DescribeCacheClustersInput{}
	result, err := session.DescribeCacheClusters(&input)
	h.LogIfError(err)
	service := "elasticache"

	var w sync.WaitGroup
	w.Add(len(result.CacheClusters))
	ch := make(chan *b.ResourceDescription, len(result.CacheClusters))
	for _, cc := range result.CacheClusters {
		go func(cc *elasticache.CacheCluster, wg *sync.WaitGroup) {
			defer wg.Done()

			resource := strings.Join([]string{"cluster", *cc.CacheClusterId}, ":")
			arn, err := nd.Parent.BuildARN(&service, &resource)
			h.LogIfError(err)

			input := elasticache.ListTagsForResourceInput{
				ResourceName: aws.String(arn),
			}
			tags, err := session.ListTagsForResource(&input)
			h.LogIfError(err)

			tl, found := nd.Parent.TagsFound(tags)
			ts := b.TagsToString(tl)

			if found {
				if r, err := createResourceDescription(nd, ts, cc); err == nil {
					ch <- r
				}
				h.LogIfError(err)
			}
		}(cc, &w)
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
