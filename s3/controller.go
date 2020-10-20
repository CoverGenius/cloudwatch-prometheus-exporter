package s3

import (
	"sync"

	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func createResourceDescription(nd *b.NamespaceDescription, tags *string, bucket *s3.Bucket) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("BucketName"),
			Value: bucket.Name,
		},
		{
			Name:  aws.String("StorageType"),
			Value: aws.String("AllStorageTypes"),
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogIfError(err)
	rd.ID = bucket.Name
	rd.Name = bucket.Name
	rd.Type = aws.String("s3")
	rd.Parent = nd
	rd.Tags = tags

	return &rd, err
}

// CreateResourceList fetches a list of all S3 buckets in the region
//
// TODO channel can be added instead of sync.WaitGroup
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	log.Debug("Creating S3 resource list ...")
	defer wg.Done()

	session := s3.New(nd.Parent.Session)
	input := s3.ListBucketsInput{}
	result, err := session.ListBuckets(&input)
	h.LogIfError(err)

	tagError := "NoSuchTagSet"

	var w sync.WaitGroup
	w.Add(len(result.Buckets))

	ch := make(chan *b.ResourceDescription, len(result.Buckets))
	for _, bucket := range result.Buckets {
		go func(bucket *s3.Bucket, wg *sync.WaitGroup) {
			defer wg.Done()
			input := s3.GetBucketLocationInput{
				Bucket: bucket.Name,
			}
			location, err := session.GetBucketLocation(&input)
			h.LogIfError(err)

			if location.LocationConstraint == nil || *location.LocationConstraint != *nd.Parent.Region {
				return
			}

			locationInput := s3.GetBucketTaggingInput{
				Bucket: bucket.Name,
			}

			tags, err := session.GetBucketTagging(&locationInput)
			if !b.IsSameErrorType(err, &tagError) {
				h.LogIfError(err)
			}

			tl, found := nd.Parent.TagsFound(tags)
			ts := b.TagsToString(tl)

			if found {
				if r, err := createResourceDescription(nd, ts, bucket); err == nil {
					ch <- r
				}
				h.LogIfError(err)
			}
		}(bucket, &w)
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
