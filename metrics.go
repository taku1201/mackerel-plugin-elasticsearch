package elasticsearch

type ClusterHealth struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               float64 `json:"number_of_nodes"`
	NumberOfDataNodes           float64 `json:"number_of_data_nodes"`
	ActivePrimaryShards         float64 `json:"active_primary_shards"`
	ActiveShards                float64 `json:"active_shards"`
	RelocatingShards            float64 `json:"relocating_shards"`
	InitializingShards          float64 `json:"initializing_shards"`
	UnassignedShards            float64 `json:"unassigned_shards"`
	DelayedUnassignedShards     float64 `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        float64 `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       float64 `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis float64 `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}
