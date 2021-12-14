package vpc

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"strconv"
	"strings"
	"time"

	"errors"
)

func findSubnet(subnets []*ec2.Subnet, id *string) (*ec2.Subnet, error) {
	for _, subnet := range subnets {
		if *subnet.SubnetId == *id {
			return subnet, nil
		}
	}

	return nil, errors.New("Subnet not found")
}

func gatherAvailableIPFunc(rds []*b.ResourceDescription, start time.Time, end time.Time) ([]*b.NonCloudWatchMetric, error) {
	if len(rds) < 1 {
		return []*b.NonCloudWatchMetric{&b.NonCloudWatchMetric{}}, nil
	}

	subnetIDs := make([]*string, len(rds))

	for idx, rd := range rds {
		subnetIDs[idx] = rd.ID
	}

	input := ec2.DescribeSubnetsInput{
		SubnetIds: subnetIDs,
	}

	session := ec2.New(rds[0].Parent.Parent.Session)
	output, err := session.DescribeSubnets(&input)
	if err != nil {
		return nil, err
	}

	result := make([]*b.NonCloudWatchMetric, len(output.Subnets))

	for idx, rd := range rds {
		subnet, err := findSubnet(output.Subnets, rd.ID)
		if err != nil {
			return result, err
		}
		metric := b.NonCloudWatchMetric{
			Timestamps: []*time.Time{aws.Time(time.Now())},
			Label: aws.String((&b.AwsLabels{
				Statistic: "Average",
				Name:      *rd.Name,
				Id:        *rd.ID,
				RType:     *rd.Type,
				Region:    *rd.Parent.Parent.Region,
				Tags:      *rd.Tags,
			}).String()),
			Values: []*float64{aws.Float64(float64(*subnet.AvailableIpAddressCount))},
		}
		result[idx] = &metric
	}

	return result, nil

}

func gatherTotalIPFunc(rds []*b.ResourceDescription, start time.Time, end time.Time) ([]*b.NonCloudWatchMetric, error) {
	if len(rds) < 1 {
		return []*b.NonCloudWatchMetric{&b.NonCloudWatchMetric{}}, nil
	}

	subnetIDs := make([]*string, len(rds))

	for idx, rd := range rds {
		subnetIDs[idx] = rd.ID
	}

	input := ec2.DescribeSubnetsInput{
		SubnetIds: subnetIDs,
	}

	session := ec2.New(rds[0].Parent.Parent.Session)
	output, err := session.DescribeSubnets(&input)
	if err != nil {
		return nil, err
	}

	result := make([]*b.NonCloudWatchMetric, len(output.Subnets))

	for idx, rd := range rds {
		subnet, err := findSubnet(output.Subnets, rd.ID)
		if err != nil {
			return result, err
		}

		shards := strings.Split(*subnet.CidrBlock, "/")
		netmask, err := strconv.Atoi(shards[1])
		if err != nil {
			return result, err
		}

		metric := b.NonCloudWatchMetric{
			Timestamps: []*time.Time{aws.Time(time.Now())},
			Label: aws.String((&b.AwsLabels{
				Statistic: "Average",
				Name:      *rd.Name,
				Id:        *rd.ID,
				RType:     *rd.Type,
				Region:    *rd.Parent.Parent.Region,
				Tags:      *rd.Tags,
			}).String()),
			Values: []*float64{aws.Float64(float64(h.IntPow(2, 32-netmask)))},
		}
		result[idx] = &metric
	}

	return result, nil
}

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"AvailableIpAddressCount": {
		Help:       aws.String("The number of ip address available for allocation."),
		OutputName: aws.String("available_ip_address_count"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.NON_CLOUDWATCH_KIND),
		GatherFunc: gatherAvailableIPFunc,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"TotalIpAddressCount": {
		Help:       aws.String("The total number of ip addresses"),
		OutputName: aws.String("total_ip_address_count"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.NON_CLOUDWATCH_KIND),
		GatherFunc: gatherTotalIPFunc,

		Dimensions: []*cloudwatch.Dimension{},
	},
}
