package ec2

import (
        b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	"github.com/aws/aws-sdk-go/aws"
)

var metrics = map[string]*b.MetricDescription{
	"CPUCreditBalance": {
		Help:       aws.String("The number of earned CPU credits that an instance has accrued since it was launched or started"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_cpu_credit_balance"),
		Data:       map[string][]*string{},
	},
	"CPUCreditUsage": {
		Help:       aws.String("The number of CPU credits spent by the instance for CPU utilization"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_cpu_credit_usage"),
		Data:       map[string][]*string{},
	},
	"CPUSurplusCreditBalance": {
		Help:       aws.String("The number of surplus credits that have been spent by an unlimited instance when its CPUCreditBalance value is zero"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_cpu_surplus_credit_balance"),
		Data:       map[string][]*string{},
	},
	"CPUSurplusCreditsCharged": {
		Help:       aws.String("The number of spent surplus credits that are not paid down by earned CPU credits, and which thus incur an additional charge"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_cpu_surplus_credits_charged"),
		Data:       map[string][]*string{},
	},
	"CPUUtilization": {
		Help:       aws.String("The percentage of allocated EC2 compute units that are currently in use on the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_cpu_utilization"),
		Data:       map[string][]*string{},
	},
	"DiskReadBytes": {
		Help:       aws.String("Bytes read from all instance store volumes available to the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_disk_read_bytes"),
		Data:       map[string][]*string{},
	},
	"DiskReadOps": {
		Help:       aws.String("Completed read operations from all instance store volumes available to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_disk_read_ops"),
		Data:       map[string][]*string{},
	},
	"DiskWriteBytes": {
		Help:       aws.String("Bytes written to all instance store volumes available to the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_disk_write_bytes"),
		Data:       map[string][]*string{},
	},
	"DiskWriteOps": {
		Help:       aws.String("Completed write operations to all instance store volumes available to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_disk_write_ops"),
		Data:       map[string][]*string{},
	},
	"EBSByteBalance": {
		Help:       aws.String("Available only for the smaller instance sizes. Provides information about the percentage of throughput credits remaining in the burst bucket"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_byte_balance"),
		Data:       map[string][]*string{},
	},
	"EBSIOBalance": {
		Help:       aws.String("Available only for the smaller instance sizes. Provides information about the percentage of I/O credits remaining in the burst bucket"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_io_balance"),
		Data:       map[string][]*string{},
	},
	"EBSReadBytes": {
		Help:       aws.String("Bytes read from all EBS volumes attached to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_read_bytes"),
		Data:       map[string][]*string{},
	},
	"EBSReadOps": {
		Help:       aws.String("Completed read operations from all Amazon EBS volumes attached to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_read_ops"),
		Data:       map[string][]*string{},
	},
	"EBSWriteBytes": {
		Help:       aws.String("Bytes written to all EBS volumes attached to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_write_bytes"),
		Data:       map[string][]*string{},
	},
	"EBSWriteOps": {
		Help:       aws.String("Completed write operations to all EBS volumes attached to the instance in a specified period of time"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_ebs_write_ops"),
		Data:       map[string][]*string{},
	},
	"NetworkIn": {
		Help:       aws.String("The number of bytes received on all network interfaces by the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_network_in"),
		Data:       map[string][]*string{},
	},
	"NetworkOut": {
		Help:       aws.String("The number of bytes sent out on all network interfaces by the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_network_out"),
		Data:       map[string][]*string{},
	},
	"NetworkPacketsIn": {
		Help:       aws.String("The number of packets received on all network interfaces by the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_network_packets_in"),
		Data:       map[string][]*string{},
	},
	"NetworkPacketsOut": {
		Help:       aws.String("The number of packets sent out on all network interfaces by the instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_network_packets_out"),
		Data:       map[string][]*string{},
	},
	"StatusCheckFailed": {
		Help:       aws.String("Reports whether the instance has passed both the instance status check and the system status check in the last minute"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_status_check_failed"),
		Data:       map[string][]*string{},
	},
	"StatusCheckFailed_Instance": {
		Help:       aws.String("Reports whether the instance has passed the instance status check in the last minute"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_status_check_failed_instance"),
		Data:       map[string][]*string{},
	},
	"StatusCheckFailed_System": {
		Help:       aws.String("Reports whether the instance has passed the system status check in the last minute"),
		Type:       aws.String("counter"),
		OutputName: aws.String("ec2_status_check_failed_system"),
		Data:       map[string][]*string{},
	},
}

func GetMetrics() map[string]*b.MetricDescription {
	return metrics
}
