package ec2

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"CPUCreditBalance": {
		Help:          aws.String("The number of earned CPU credits that an instance has accrued since it was launched or started"),
		OutputName:    aws.String("ec2_cpu_credit_balance"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUCreditUsage": {
		Help:          aws.String("The number of CPU credits spent by the instance for CPU utilization"),
		OutputName:    aws.String("ec2_cpu_credit_usage"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditBalance": {
		Help:          aws.String("The number of surplus credits that have been spent by an unlimited instance when its CPUCreditBalance value is zero"),
		OutputName:    aws.String("ec2_cpu_surplus_credit_balance"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditsCharged": {
		Help:          aws.String("The number of spent surplus credits that are not paid down by earned CPU credits, and which thus incur an additional charge"),
		OutputName:    aws.String("ec2_cpu_surplus_credits_charged"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUUtilization": {
		Help:          aws.String("The percentage of allocated EC2 compute units that are currently in use on the instance"),
		OutputName:    aws.String("ec2_cpu_utilization"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskReadBytes": {
		Help:          aws.String("Bytes read from all instance store volumes available to the instance"),
		OutputName:    aws.String("ec2_disk_read_bytes"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskReadOps": {
		Help:          aws.String("Completed read operations from all instance store volumes available to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_disk_read_ops"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskWriteBytes": {
		Help:          aws.String("Bytes written to all instance store volumes available to the instance"),
		OutputName:    aws.String("ec2_disk_write_bytes"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskWriteOps": {
		Help:          aws.String("Completed write operations to all instance store volumes available to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_disk_write_ops"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSByteBalance": {
		Help:          aws.String("Available only for the smaller instance sizes. Provides information about the percentage of throughput credits remaining in the burst bucket"),
		OutputName:    aws.String("ec2_ebs_byte_balance"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSIOBalance": {
		Help:          aws.String("Available only for the smaller instance sizes. Provides information about the percentage of I/O credits remaining in the burst bucket"),
		OutputName:    aws.String("ec2_ebs_io_balance"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSReadBytes": {
		Help:          aws.String("Bytes read from all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_ebs_read_bytes"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSReadOps": {
		Help:          aws.String("Completed read operations from all Amazon EBS volumes attached to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_ebs_read_ops"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSWriteBytes": {
		Help:          aws.String("Bytes written to all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_ebs_write_bytes"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSWriteOps": {
		Help:          aws.String("Completed write operations to all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    aws.String("ec2_ebs_write_ops"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkIn": {
		Help:          aws.String("The number of bytes received on all network interfaces by the instance"),
		OutputName:    aws.String("ec2_network_in"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkOut": {
		Help:          aws.String("The number of bytes sent out on all network interfaces by the instance"),
		OutputName:    aws.String("ec2_network_out"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkPacketsIn": {
		Help:          aws.String("The number of packets received on all network interfaces by the instance"),
		OutputName:    aws.String("ec2_network_packets_in"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkPacketsOut": {
		Help:          aws.String("The number of packets sent out on all network interfaces by the instance"),
		OutputName:    aws.String("ec2_network_packets_out"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed": {
		Help:          aws.String("Reports whether the instance has passed both the instance status check and the system status check in the last minute"),
		OutputName:    aws.String("ec2_status_check_failed"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed_Instance": {
		Help:          aws.String("Reports whether the instance has passed the instance status check in the last minute"),
		OutputName:    aws.String("ec2_status_check_failed_instance"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed_System": {
		Help:          aws.String("Reports whether the instance has passed the system status check in the last minute"),
		OutputName:    aws.String("ec2_status_check_failed_system"),
		Statistic:     h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
}
