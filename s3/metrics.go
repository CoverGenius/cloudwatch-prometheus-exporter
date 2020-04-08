package s3

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var metrics = map[string]*b.MetricDescription{
	"BucketSizeBytes": {
		Help:       aws.String("The amount of data in bytes stored in a bucket in the STANDARD storage class, INTELLIGENT_TIERING storage class, Standard - Infrequent Access (STANDARD_IA) storage class, OneZone - Infrequent Access (ONEZONE_IA), Reduced Redundancy Storage (RRS) class, Deep Archive Storage (DEEP_ARCHIVE) class or, Glacier (GLACIER) storage class"),
		OutputName: aws.String("s3_bucket_size_bytes"),
		Statistic:  h.StringPointers("Average"),
		Period:     1440,
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			},
		},
	},
	"NumberOfObjects": {
		Help:       aws.String("The total number of objects stored in a bucket for all storage classes except for the GLACIER storage class"),
		OutputName: aws.String("s3_number_of_objects"),
		Statistic:  h.StringPointers("Average"),
		Period:     1440,
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("AllStorageTypes"),
			},
		},
	},
}

// GetMetrics returns a map of MetricDescriptions to be exported for this namespace
func GetMetrics() map[string]*b.MetricDescription {
	return metrics
}
