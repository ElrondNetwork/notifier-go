package common

const (
	// WSPublisherType defines a webserver api type using WebSockets
	WSPublisherType string = "ws"

	// MessageQueuePublisherType defines a webserver api type using a message queueing service
	MessageQueuePublisherType string = "rabbitmq"

	// WSPublisherTypeV2 defines a webserver api type using WebSockets without custom filtering and subscriptions
	WSPublisherTypeV2 string = "wsv2"
)

const (
	// RedisInstanceConnType specifies a redis connection to a single instance
	RedisInstanceConnType string = "instance"

	// RedisSentinelConnType specifies a redis connection to a setup with sentinel
	RedisSentinelConnType string = "sentinel"
)

const (
	// PushLogsAndEvents defines the subscription event type for pushing block events
	PushLogsAndEvents string = "all_events"

	// BlockEvents defines the subscription event type for block info with logs and events
	BlockEvents string = "block_events"

	// RevertBlockEvents defines the subscription event type for revert block
	RevertBlockEvents string = "revert_events"

	// FinalizedBlockEvents defines the subscription event type for finalized blocks
	FinalizedBlockEvents string = "finalized_events"

	// BlockTxs defines the subscription event type for block txs
	BlockTxs string = "block_txs"

	// BlockScrs defines the subscription event type for block scrs
	BlockScrs string = "block_scrs"
)

const (
	// WSObsConnectorType defines the websocket observer connector type
	WSObsConnectorType string = "ws"

	// HTTPConnectorType defines the http observer connector type
	HTTPConnectorType string = "http"
)

const (
	// PayloadV0 defines the version of payload before versioning implementation
	PayloadV0 uint32 = 0

	// PayloadV1 defines first payload implementation with versioning
	PayloadV1 uint32 = 1
)
