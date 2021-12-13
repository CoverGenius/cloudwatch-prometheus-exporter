package rds

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"BinLogDiskUsage": {
		Help:       aws.String("The amount of disk space occupied by binary logs on the master. Applies to MySQL read replicas"),
		OutputName: aws.String("rds_bin_log_disk_usage"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"BurstBalance": {
		Help:       aws.String("The percent of General Purpose SSD (gp2) burst-bucket I/O credits available"),
		OutputName: aws.String("rds_burst_balance"),
		Statistic:  h.StringPointers("Average", "Minimum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUCreditBalance": {
		Help:       aws.String("The number of earned CPU credits that an instance has accrued. This represents the number of credits currently available."),
		OutputName: aws.String("rds_cpu_credit_balance"),
		Statistic:  h.StringPointers("Average", "Minimum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUCreditUsage": {
		Help:       aws.String("The number of CPU credits spent by the instance for CPU utilization. One CPU credit equals one vCPU running at 100 percent utilization for one minute or an equivalent combination of vCPUs, utilization, and time"),
		OutputName: aws.String("rds_cpu_credit_usage"),
		Statistic:  h.StringPointers("Average", "Sum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditBalance": {
		Help:       aws.String("The number of surplus credits that have been spent by an unlimited instance when its CPUCreditBalance value is zero"),
		OutputName: aws.String("rds_cpu_surplus_credit_balance"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUSurplusCreditsCharged": {
		Help:       aws.String("The number of spent surplus credits that are not paid down by earned CPU credits, and which thus incur an additional charge"),
		OutputName: aws.String("rds_cpu_surplus_credits_charged"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUUtilization": {
		Help:       aws.String("The percentage of CPU utilization"),
		OutputName: aws.String("rds_cpu_utilization"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DatabaseConnections": {
		Help:       aws.String("The number of database connections in use"),
		OutputName: aws.String("rds_database_connections"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DBLoad": {
		Help:       aws.String("The number of active sessions for the DB engine. Typically, you want the data for the average number of active sessions"),
		OutputName: aws.String("rds_db_load"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DBLoadCPU": {
		Help:       aws.String("The number of active sessions where the wait event type is CPU"),
		OutputName: aws.String("rds_db_load_cpu"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DBLoadNonCPU": {
		Help:       aws.String("The number of active sessions where the wait event type is not CPU"),
		OutputName: aws.String("rds_db_load_non_cpu"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DiskQueueDepth": {
		Help:       aws.String("The number of outstanding IOs (read/write requests) waiting to access the disk"),
		OutputName: aws.String("rds_disk_queue_depth"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"FreeableMemory": {
		Help:       aws.String("The amount of available random access memory"),
		OutputName: aws.String("rds_freeable_memory"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"FreeStorageSpace": {
		Help:       aws.String("The amount of available storage space"),
		OutputName: aws.String("rds_free_storage_space"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"MaximumUsedTransactionIDs": {
		Help:       aws.String("The maximum transaction ID that has been used. Applies to PostgreSQL"),
		OutputName: aws.String("rds_maximum_used_transaction_ids"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkReceiveThroughput": {
		Help:       aws.String("The incoming (Receive) network traffic on the DB instance, including both customer database traffic and Amazon RDS traffic used for monitoring and replication"),
		OutputName: aws.String("rds_network_receive_throughput"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkTransmitThroughput": {
		Help:       aws.String("The outgoing (Transmit) network traffic on the DB instance, including both customer database traffic and Amazon RDS traffic used for monitoring and replication"),
		OutputName: aws.String("rds_network_transmit_throughput"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"OldestReplicationSlotLag": {
		Help:       aws.String("The lagging size of the replica lagging the most in terms of WAL data received. Applies to PostgreSQL"),
		OutputName: aws.String("rds_oldest_replication_slot_lag"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReadIOPS": {
		Help:       aws.String("The average number of disk read I/O operations per second"),
		OutputName: aws.String("rds_read_iops"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReadLatency": {
		Help:       aws.String("The amount of time taken per disk I/O operation"),
		OutputName: aws.String("rds_read_latency"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReadThroughput": {
		Help:       aws.String("The number of bytes read from disk per second"),
		OutputName: aws.String("rds_read_throughput"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReplicaLag": {
		Help:       aws.String("The amount of time a Read Replica DB instance lags behind the source DB instance. Applies to MySQL, MariaDB, and PostgreSQL Read Replicas"),
		OutputName: aws.String("rds_replica_lag"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReplicationSlotDiskUsage": {
		Help:       aws.String("The disk space used by replication slot files. Applies to PostgreSQL"),
		OutputName: aws.String("rds_replication_slot_disk_usage"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SwapUsage": {
		Help:       aws.String("The amount of swap space used on the DB instance. This metric is not available for SQL Server"),
		OutputName: aws.String("rds_swap_usage"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"TransactionLogsDiskUsage": {
		Help:       aws.String("The disk space used by transaction logs. Applies to PostgreSQL"),
		OutputName: aws.String("rds_transaction_logs_disk_usage"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"TransactionLogsGeneration": {
		Help:       aws.String("The size of transaction logs generated per second. Applies to PostgreSQL"),
		OutputName: aws.String("rds_transaction_logs_generation"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"WriteIOPS": {
		Help:       aws.String("The average number of disk write I/O operations per second"),
		OutputName: aws.String("rds_write_iops"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"WriteLatency": {
		Help:       aws.String("The amount of time taken per disk I/O operation"),
		OutputName: aws.String("rds_write_latency"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"WriteThroughput": {
		Help:       aws.String("The number of bytes written to disk per second"),
		OutputName: aws.String("rds_write_throughput"),
		Statistic:  h.StringPointers("Average", "Maximum"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
}
