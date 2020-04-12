package s3

import (
	"sync"

	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func CreateResourceDescription(nd *b.NamespaceDescription, bucket *s3.Bucket) error {
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
	h.LogError(err)
	rd.ID = bucket.Name
	rd.Name = bucket.Name
	rd.Type = aws.String("s3")
	rd.Parent = nd
	nd.Mutex.Lock()
	nd.Resources = append(nd.Resources, &rd)
	nd.Mutex.Unlock()

	return nil
}

// CreateResourceList fetches a list of all S3 buckets in the region
//
// TODO channel can be added instead of sync.WaitGroup
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) error {
	log.Debug("Creating S3 resource list ...")
	defer wg.Done()

	nd.Resources = []*b.ResourceDescription{}
	session := s3.New(nd.Parent.Session)
	input := s3.ListBucketsInput{}
	result, err := session.ListBuckets(&input)
	h.LogError(err)

	tagError := "NoSuchTagSet"

	var w sync.WaitGroup
	w.Add(len(result.Buckets))

	for _, bucket := range result.Buckets {
		go func(bucket *s3.Bucket, wg *sync.WaitGroup) {
			defer wg.Done()
			input := s3.GetBucketLocationInput{
				Bucket: bucket.Name,
			}
			location, err := session.GetBucketLocation(&input)
			h.LogError(err)

			if location.LocationConstraint == nil || *location.LocationConstraint != *nd.Parent.Region {
				return
			}

			locationInput := s3.GetBucketTaggingInput{
				Bucket: bucket.Name,
			}

			tags, err := session.GetBucketTagging(&locationInput)
			if b.IsSameErrorType(err, &tagError) == false {
				h.LogError(err)
			}

			if nd.Parent.TagsFound(tags) {
				err := CreateResourceDescription(nd, bucket)
				h.LogError(err)
			}
		}(bucket, &w)
	}
	w.Wait()

	return nil
}
