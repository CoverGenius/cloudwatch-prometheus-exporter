package elasticache

import (
        b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	"github.com/aws/aws-sdk-go/aws"
)

var metrics = map[string]*b.MetricDescription{
	"ActiveDefragHits": {
		Help:       aws.String("The number of value reallocations per minute performed by the active defragmentation process"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_active_defrag_hits"),
		Data:       map[string][]*string{},
	},
	"BytesUsedForCache": {
		Help:       aws.String("The total number of bytes allocated by Redis for all purposes, including the dataset, buffers, etc"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_bytes_used_for_cache"),
		Data:       map[string][]*string{},
	},
	"CacheHits": {
		Help:       aws.String("The number of successful read-only key lookups in the main dictionary"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_cache_hits"),
		Data:       map[string][]*string{},
	},
	"CacheMisses": {
		Help:       aws.String("The number of unsuccessful read-only key lookups in the main dictionary"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_cache_misses"),
		Data:       map[string][]*string{},
	},
	"CPUUtilization": {
		Help:       aws.String("The percentage of CPU utilization"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_cpu_utilization"),
		Data:       map[string][]*string{},
	},
	"CurrConnections": {
		Help:       aws.String("The number of client connections, excluding connections from read replicas. ElastiCache uses two to three of the connections to monitor the cluster in each case"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_curr_connections"),
		Data:       map[string][]*string{},
	},
	"CurrItems": {
		Help:       aws.String("The number of items in the cache"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_curr_items"),
		Data:       map[string][]*string{},
	},
	"EngineCPUUtilization": {
		Help:       aws.String("Provides CPU utilization of the Redis engine thread. Since Redis is single-threaded, you can use this metric to analyze the load of the Redis process itself"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_engine_cpu_utilization"),
		Data:       map[string][]*string{},
	},
	"Evictions": {
		Help:       aws.String("The number of keys that have been evicted due to the maxmemory limit"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_evictions"),
		Data:       map[string][]*string{},
	},
	"FreeableMemory": {
		Help:       aws.String("The amount of free memory available on the host"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_freeable_memory"),
		Data:       map[string][]*string{},
	},
	"GetTypeCmds": {
		Help:       aws.String("The total number of read-only type commands"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_get_type_cmds"),
		Data:       map[string][]*string{},
	},
	"IsMaster": {
		Help:       aws.String("Returns 1 in case if node is master"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_is_master"),
		Data:       map[string][]*string{},
	},
	"KeyBasedCmds": {
		Help:       aws.String("The total number of commands that are key-based"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_key_based_cmds"),
		Data:       map[string][]*string{},
	},
	"ListBasedCmds": {
		Help:       aws.String("The total number of commands that are list-based"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_list_based_cmds"),
		Data:       map[string][]*string{},
	},
	"MasterLinkHealthStatus": {
		Help:       aws.String("This status has two values: 0 or 1. The value 0 indicates that data in the Elasticache primary node is not in sync with Redis on EC2. The value of 1 indicates that the data is in sync"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_master_link_health_status"),
		Data:       map[string][]*string{},
	},
	"NetworkBytesIn": {
		Help:       aws.String("The number of bytes the host has read from the network"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_network_bytes_in"),
		Data:       map[string][]*string{},
	},
	"NetworkBytesOut": {
		Help:       aws.String("The number of bytes the host has written to the network"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_network_bytes_out"),
		Data:       map[string][]*string{},
	},
	"NetworkPacketsIn": {
		Help:       aws.String("The number of packets received on all network interfaces by the instance. This metric identifies the volume of incoming traffic in terms of the number of packets on a single instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_network_packets_in"),
		Data:       map[string][]*string{},
	},
	"NetworkPacketsOut": {
		Help:       aws.String("The number of packets sent out on all network interfaces by the instance. This metric identifies the volume of outgoing traffic in terms of the number of packets on a single instance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_network_packets_out"),
		Data:       map[string][]*string{},
	},
	"NewConnections": {
		Help:       aws.String("The total number of connections that have been accepted by the server during this period"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_new_connections"),
		Data:       map[string][]*string{},
	},
	"Reclaimed": {
		Help:       aws.String("The total number of key expiration events"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_reclaimed"),
		Data:       map[string][]*string{},
	},
	"ReplicationBytes": {
		Help:       aws.String("For nodes in a replicated configuration, ReplicationBytes reports the number of bytes that the primary is sending to all of its replicas. This metric is representative of the write load on the replication group"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_replication_bytes"),
		Data:       map[string][]*string{},
	},
	"ReplicationLag": {
		Help:       aws.String("This metric is only applicable for a node running as a read replica. It represents how far behind, in seconds, the replica is in applying changes from the primary node"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_replication_lag"),
		Data:       map[string][]*string{},
	},
	"SaveInProgress": {
		Help:       aws.String("This binary metric returns 1 whenever a background save (forked or forkless) is in progress, and 0 otherwise. A background save process is typically used during snapshots and syncs. These operations can cause degraded performance"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_save_in_progress"),
		Data:       map[string][]*string{},
	},
	"SetBasedCmds": {
		Help:       aws.String("The total number of commands that are set-based"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_set_based_cmds"),
		Data:       map[string][]*string{},
	},
	"SetTypeCmds": {
		Help:       aws.String("The total number of write types of commands"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_set_type_cmds"),
		Data:       map[string][]*string{},
	},
	"SortedSetBasedCmds": {
		Help:       aws.String("The total number of commands that are sorted set-based"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_sorted_set_based_cmds"),
		Data:       map[string][]*string{},
	},
	"StringBasedCmds": {
		Help:       aws.String("The total number of commands that are string-based"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_string_based_cmds"),
		Data:       map[string][]*string{},
	},
	"SwapUsage": {
		Help:       aws.String("The amount of swap used on the host"),
		Type:       aws.String("counter"),
		OutputName: aws.String("elasticache_swap_usage"),
		Data:       map[string][]*string{},
	},
}

func GetMetrics() map[string]*b.MetricDescription {
	return metrics
}
