// package codes defines the error types dynago maps.
// This is an internal package to allow fast iteration on the abstraction
// without having to commit to the interface.
package codes

type ErrorCode int

const (
	ErrorUnknown ErrorCode = iota

	ErrorConditionFailed        // Conditional put/update failed; condition not met
	ErrorCollectionSizeExceeded // Item collection (local secondary index) too large
	ErrorThroughputExceeded     // Exceeded provisioned throughput for table or shard
	ErrorNotFound               // Resource referenced by key not found
	ErrorInternalFailure        // Internal server error
	ErrorAuth                   // Encapsulates various authorization errors
	ErrorInvalidParameter       // Encapsulates many forms of invalid input errors
	ErrorServiceUnavailable     // Amazon service unavailable
	ErrorThrottling             // Amazon is throttling us, try later
	ErrorResourceInUse          // Tried to create existing table, delete a table in CREATING state, etc.

	// DynamoDB Streams
	ErrorExpiredIterator // Iterator is no longer valid
	ErrorTrimmedData     // Attempted to access data older than 24h
)
