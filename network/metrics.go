package network

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.ConfigMetric{
	"ActiveConnectionCount": {
		Help:          ("The total number of concurrent active TCP connections through the NAT gateway"),
		OutputName:    ("nat_gateway_active_connection_count"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"BytesInFromDestination": {
		Help:          ("The number of bytes received by the NAT gateway from the destination"),
		OutputName:    ("nat_gateway_bytes_in_from_destination"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"BytesInFromSource": {
		Help:          ("The number of bytes received by the NAT gateway from clients in your VPC"),
		OutputName:    ("nat_gateway_bytes_in_from_source"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"BytesOutToDestination": {
		Help:          ("The number of bytes sent out through the NAT gateway to the destination"),
		OutputName:    ("nat_gateway_bytes_out_to_destination"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"BytesOutToSource": {
		Help:          ("The number of bytes sent through the NAT gateway to the clients in your VPC"),
		OutputName:    ("nat_gateway_bytes_out_to_source"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"ConnectionAttemptCount": {
		Help:          ("The number of connection attempts made through the NAT gateway"),
		OutputName:    ("nat_gateway_connection_attempt_count"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"ConnectionEstablishedCount": {
		Help:          ("The number of connections established through the NAT gateway"),
		OutputName:    ("nat_gateway_connection_established_count"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"ErrorPortAllocation": {
		Help:          ("The number of times the NAT gateway could not allocate a source port"),
		OutputName:    ("nat_gateway_error_port_allocation"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"IdleTimeoutCount": {
		Help:          ("The number of connections that transitioned from the active state to the idle state. An active connection transitions to idle if it was not closed gracefully and there was no activity for the last 350 seconds"),
		OutputName:    ("nat_gateway_idle_timeout_count"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"PacketsDropCount": {
		Help:          ("The number of packets dropped by the NAT gateway"),
		OutputName:    ("nat_gateway_packets_drop_count"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"PacketsInFromDestination": {
		Help:          ("The number of packets received by the NAT gateway from the destination"),
		OutputName:    ("nat_gateway_packets_in_from_destination"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"PacketsInFromSource": {
		Help:          ("The number of packets received by the NAT gateway from clients in your VPC"),
		OutputName:    ("nat_gateway_packets_in_from_source"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"PacketsOutToDestination": {
		Help:          ("The number of packets sent out through the NAT gateway to the destination"),
		OutputName:    ("nat_gateway_packets_out_to_destination"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"PacketsOutToSource": {
		Help:          ("The number of packets sent through the NAT gateway to the clients in your VPC"),
		OutputName:    ("nat_gateway_packets_out_to_source"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
}
