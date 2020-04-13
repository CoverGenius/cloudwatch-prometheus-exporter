package s3

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.ConfigMetric{
	"BucketSizeBytes": {
		Help:          "The amount of data in bytes stored in a bucket in the STANDARD storage class, INTELLIGENT_TIERING storage class, Standard - Infrequent Access (STANDARD_IA) storage class, OneZone - Infrequent Access (ONEZONE_IA), Reduced Redundancy Storage (RRS) class, Deep Archive Storage (DEEP_ARCHIVE) class or, Glacier (GLACIER) storage class",
		OutputName:    "s3_bucket_size_bytes",
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 60 * 60 * 24,
		RangeSeconds:  60 * 60 * 24 * 7,
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			},
		},
	},
	"NumberOfObjects": {
		Help:          ("The total number of objects stored in a bucket for all storage classes except for the GLACIER storage class"),
		OutputName:    ("s3_number_of_objects"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 60 * 60 * 24,
		RangeSeconds:  60 * 60 * 24 * 7,
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("AllStorageTypes"),
			},
		},
	},
}
