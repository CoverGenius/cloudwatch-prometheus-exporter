package network

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var metrics = map[string]*b.MetricDescription{
	"ActiveConnectionCount": {
		Help:       aws.String("The total number of concurrent active TCP connections through the NAT gateway"),
		OutputName: aws.String("nat_gateway_active_connection_count"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"BytesInFromDestination": {
		Help:       aws.String("The number of bytes received by the NAT gateway from the destination"),
		OutputName: aws.String("nat_gateway_bytes_in_from_destination"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"BytesInFromSource": {
		Help:       aws.String("The number of bytes received by the NAT gateway from clients in your VPC"),
		OutputName: aws.String("nat_gateway_bytes_in_from_source"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"BytesOutToDestination": {
		Help:       aws.String("The number of bytes sent out through the NAT gateway to the destination"),
		OutputName: aws.String("nat_gateway_bytes_out_to_destination"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"BytesOutToSource": {
		Help:       aws.String("The number of bytes sent through the NAT gateway to the clients in your VPC"),
		OutputName: aws.String("nat_gateway_bytes_out_to_source"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"ConnectionAttemptCount": {
		Help:       aws.String("The number of connection attempts made through the NAT gateway"),
		OutputName: aws.String("nat_gateway_connection_attempt_count"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"ConnectionEstablishedCount": {
		Help:       aws.String("The number of connections established through the NAT gateway"),
		OutputName: aws.String("nat_gateway_connection_established_count"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"ErrorPortAllocation": {
		Help:       aws.String("The number of times the NAT gateway could not allocate a source port"),
		OutputName: aws.String("nat_gateway_error_port_allocation"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"IdleTimeoutCount": {
		Help:       aws.String("The number of connections that transitioned from the active state to the idle state. An active connection transitions to idle if it was not closed gracefully and there was no activity for the last 350 seconds"),
		OutputName: aws.String("nat_gateway_idle_timeout_count"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"PacketsDropCount": {
		Help:       aws.String("The number of packets dropped by the NAT gateway"),
		OutputName: aws.String("nat_gateway_packets_drop_count"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"PacketsInFromDestination": {
		Help:       aws.String("The number of packets received by the NAT gateway from the destination"),
		OutputName: aws.String("nat_gateway_packets_in_from_destination"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"PacketsInFromSource": {
		Help:       aws.String("The number of packets received by the NAT gateway from clients in your VPC"),
		OutputName: aws.String("nat_gateway_packets_in_from_source"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"PacketsOutToDestination": {
		Help:       aws.String("The number of packets sent out through the NAT gateway to the destination"),
		OutputName: aws.String("nat_gateway_packets_out_to_destination"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
	"PacketsOutToSource": {
		Help:       aws.String("The number of packets sent through the NAT gateway to the clients in your VPC"),
		OutputName: aws.String("nat_gateway_packets_out_to_source"),
		Statistic:  h.StringPointers("Average"),
		Period:     5,
		Dimensions: []*cloudwatch.Dimension{},
	},
}

// GetMetrics returns a map of MetricDescriptions to be exported for this namespace
func GetMetrics() map[string]*b.MetricDescription {
	return metrics
}
