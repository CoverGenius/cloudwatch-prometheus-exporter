package ec2

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default ConfigMetrics for this namespace
var Metrics = map[string]*b.ConfigMetric{
	"CPUCreditBalance": {
		Help:          ("The number of earned CPU credits that an instance has accrued since it was launched or started"),
		OutputName:    ("ec2_cpu_credit_balance"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUCreditUsage": {
		Help:          ("The number of CPU credits spent by the instance for CPU utilization"),
		OutputName:    ("ec2_cpu_credit_usage"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditBalance": {
		Help:          ("The number of surplus credits that have been spent by an unlimited instance when its CPUCreditBalance value is zero"),
		OutputName:    ("ec2_cpu_surplus_credit_balance"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditsCharged": {
		Help:          ("The number of spent surplus credits that are not paid down by earned CPU credits, and which thus incur an additional charge"),
		OutputName:    ("ec2_cpu_surplus_credits_charged"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"CPUUtilization": {
		Help:          ("The percentage of allocated EC2 compute units that are currently in use on the instance"),
		OutputName:    ("ec2_cpu_utilization"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskReadBytes": {
		Help:          ("Bytes read from all instance store volumes available to the instance"),
		OutputName:    ("ec2_disk_read_bytes"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskReadOps": {
		Help:          ("Completed read operations from all instance store volumes available to the instance in a specified period of time"),
		OutputName:    ("ec2_disk_read_ops"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskWriteBytes": {
		Help:          ("Bytes written to all instance store volumes available to the instance"),
		OutputName:    ("ec2_disk_write_bytes"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"DiskWriteOps": {
		Help:          ("Completed write operations to all instance store volumes available to the instance in a specified period of time"),
		OutputName:    ("ec2_disk_write_ops"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSByteBalance": {
		Help:          ("Available only for the smaller instance sizes. Provides information about the percentage of throughput credits remaining in the burst bucket"),
		OutputName:    ("ec2_ebs_byte_balance"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSIOBalance": {
		Help:          ("Available only for the smaller instance sizes. Provides information about the percentage of I/O credits remaining in the burst bucket"),
		OutputName:    ("ec2_ebs_io_balance"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSReadBytes": {
		Help:          ("Bytes read from all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    ("ec2_ebs_read_bytes"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSReadOps": {
		Help:          ("Completed read operations from all Amazon EBS volumes attached to the instance in a specified period of time"),
		OutputName:    ("ec2_ebs_read_ops"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSWriteBytes": {
		Help:          ("Bytes written to all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    ("ec2_ebs_write_bytes"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"EBSWriteOps": {
		Help:          ("Completed write operations to all EBS volumes attached to the instance in a specified period of time"),
		OutputName:    ("ec2_ebs_write_ops"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkIn": {
		Help:          ("The number of bytes received on all network interfaces by the instance"),
		OutputName:    ("ec2_network_in"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkOut": {
		Help:          ("The number of bytes sent out on all network interfaces by the instance"),
		OutputName:    ("ec2_network_out"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkPacketsIn": {
		Help:          ("The number of packets received on all network interfaces by the instance"),
		OutputName:    ("ec2_network_packets_in"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"NetworkPacketsOut": {
		Help:          ("The number of packets sent out on all network interfaces by the instance"),
		OutputName:    ("ec2_network_packets_out"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed": {
		Help:          ("Reports whether the instance has passed both the instance status check and the system status check in the last minute"),
		OutputName:    ("ec2_status_check_failed"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed_Instance": {
		Help:          ("Reports whether the instance has passed the instance status check in the last minute"),
		OutputName:    ("ec2_status_check_failed_instance"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
	"StatusCheckFailed_System": {
		Help:          ("Reports whether the instance has passed the system status check in the last minute"),
		OutputName:    ("ec2_status_check_failed_system"),
		Statistics:    h.StringPointers("Average"),
		PeriodSeconds: 300,
		Dimensions:    []*cloudwatch.Dimension{},
	},
}
