package elasticache

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"ActiveDefragHits": {
		Help:       aws.String("The number of value reallocations per minute performed by the active defragmentation process"),
		OutputName: aws.String("elasticache_active_defrag_hits"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"BytesUsedForCache": {
		Help:       aws.String("The total number of bytes allocated by Redis for all purposes, including the dataset, buffers, etc"),
		OutputName: aws.String("elasticache_bytes_used_for_cache"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CacheHits": {
		Help:       aws.String("The number of successful read-only key lookups in the main dictionary"),
		OutputName: aws.String("elasticache_cache_hits"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CacheMisses": {
		Help:       aws.String("The number of unsuccessful read-only key lookups in the main dictionary"),
		OutputName: aws.String("elasticache_cache_misses"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CPUUtilization": {
		Help:       aws.String("The percentage of CPU utilization"),
		OutputName: aws.String("elasticache_cpu_utilization"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CurrConnections": {
		Help:       aws.String("The number of client connections, excluding connections from read replicas. ElastiCache uses two to three of the connections to monitor the cluster in each case"),
		OutputName: aws.String("elasticache_curr_connections"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"CurrItems": {
		Help:       aws.String("The number of items in the cache"),
		OutputName: aws.String("elasticache_curr_items"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"EngineCPUUtilization": {
		Help:       aws.String("Provides CPU utilization of the Redis engine thread. Since Redis is single-threaded, you can use this metric to analyze the load of the Redis process itself"),
		OutputName: aws.String("elasticache_engine_cpu_utilization"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"Evictions": {
		Help:       aws.String("The number of keys that have been evicted due to the maxmemory limit"),
		OutputName: aws.String("elasticache_evictions"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"FreeableMemory": {
		Help:       aws.String("The amount of free memory available on the host"),
		OutputName: aws.String("elasticache_freeable_memory"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"DatabaseMemoryUsagePercentage": {
		Help:       aws.String("The percentage of available memory used by the database"),
		OutputName: aws.String("elasticache_database_memory_usage_percentage"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"GetTypeCmds": {
		Help:       aws.String("The total number of read-only type commands"),
		OutputName: aws.String("elasticache_get_type_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"IsMaster": {
		Help:       aws.String("Returns 1 in case if node is master"),
		OutputName: aws.String("elasticache_is_master"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"KeyBasedCmds": {
		Help:       aws.String("The total number of commands that are key-based"),
		OutputName: aws.String("elasticache_key_based_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ListBasedCmds": {
		Help:       aws.String("The total number of commands that are list-based"),
		OutputName: aws.String("elasticache_list_based_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"MasterLinkHealthStatus": {
		Help:       aws.String("This status has two values: 0 or 1. The value 0 indicates that data in the Elasticache primary node is not in sync with Redis on EC2. The value of 1 indicates that the data is in sync"),
		OutputName: aws.String("elasticache_master_link_health_status"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkBytesIn": {
		Help:       aws.String("The number of bytes the host has read from the network"),
		OutputName: aws.String("elasticache_network_bytes_in"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkBytesOut": {
		Help:       aws.String("The number of bytes the host has written to the network"),
		OutputName: aws.String("elasticache_network_bytes_out"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkPacketsIn": {
		Help:       aws.String("The number of packets received on all network interfaces by the instance. This metric identifies the volume of incoming traffic in terms of the number of packets on a single instance"),
		OutputName: aws.String("elasticache_network_packets_in"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NetworkPacketsOut": {
		Help:       aws.String("The number of packets sent out on all network interfaces by the instance. This metric identifies the volume of outgoing traffic in terms of the number of packets on a single instance"),
		OutputName: aws.String("elasticache_network_packets_out"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NewConnections": {
		Help:       aws.String("The total number of connections that have been accepted by the server during this period"),
		OutputName: aws.String("elasticache_new_connections"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"Reclaimed": {
		Help:       aws.String("The total number of key expiration events"),
		OutputName: aws.String("elasticache_reclaimed"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReplicationBytes": {
		Help:       aws.String("For nodes in a replicated configuration, ReplicationBytes reports the number of bytes that the primary is sending to all of its replicas. This metric is representative of the write load on the replication group"),
		OutputName: aws.String("elasticache_replication_bytes"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ReplicationLag": {
		Help:       aws.String("This metric is only applicable for a node running as a read replica. It represents how far behind, in seconds, the replica is in applying changes from the primary node"),
		OutputName: aws.String("elasticache_replication_lag"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SaveInProgress": {
		Help:       aws.String("This binary metric returns 1 whenever a background save (forked or forkless) is in progress, and 0 otherwise. A background save process is typically used during snapshots and syncs. These operations can cause degraded performance"),
		OutputName: aws.String("elasticache_save_in_progress"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SetBasedCmds": {
		Help:       aws.String("The total number of commands that are set-based"),
		OutputName: aws.String("elasticache_set_based_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SetTypeCmds": {
		Help:       aws.String("The total number of write types of commands"),
		OutputName: aws.String("elasticache_set_type_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SortedSetBasedCmds": {
		Help:       aws.String("The total number of commands that are sorted set-based"),
		OutputName: aws.String("elasticache_sorted_set_based_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"StringBasedCmds": {
		Help:       aws.String("The total number of commands that are string-based"),
		OutputName: aws.String("elasticache_string_based_cmds"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SwapUsage": {
		Help:       aws.String("The amount of swap used on the host"),
		OutputName: aws.String("elasticache_swap_usage"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
}
