package schema

type KeyType string

const (
	HashKey  KeyType = "HASH"
	RangeKey KeyType = "RANGE"
)

type AttributeType string

const (
	String AttributeType = "S"
	Number AttributeType = "N"
	Binary AttributeType = "B"
)

type ProjectionType string

const (
	ProjectKeysOnly ProjectionType = "KEYS_ONLY"
	ProjectInclude  ProjectionType = "INCLUDE"
	ProjectAll      ProjectionType = "ALL"
)

type TableDescription struct {
	TableName        string
	TableSizeBytes   uint64
	TableStatus      string
	CreationDateTime float64 // Unix epoch time

	KeySchema              []KeySchema
	AttributeDefinitions   []AttributeDefinition
	GlobalSecondaryIndexes []SecondaryIndexResponse
	LocalSecondaryIndexes  []SecondaryIndexResponse
	ProvisionedThroughput  ProvisionedThroughputDescription

	// Streams
	LatestStreamArn     string
	LatestStreamLabel   string
	StreamSpecification *StreamSpecification
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  uint
	WriteCapacityUnits uint
}

type AttributeDefinition struct {
	AttributeName string
	AttributeType AttributeType
}

type KeySchema struct {
	AttributeName string
	KeyType       KeyType
}

type SecondaryIndex struct {
	IndexName             string
	KeySchema             []KeySchema
	Projection            Projection
	ProvisionedThroughput ProvisionedThroughput
}

type Projection struct {
	ProjectionType   ProjectionType
	NonKeyAttributes []string `json:",omitempty"`
}

// Secondary indexes as described in table descriptions
type SecondaryIndexResponse struct {
	SecondaryIndex
	Backfilling    bool
	IndexSizeBytes int
	IndexStatus    string
	ItemCount      int
}

type ProvisionedThroughputDescription struct {
	ProvisionedThroughput
	LastDecreaseDateTime   float64
	LastIncreaseDateTime   float64
	NumberOfDecreasesToday int
}

type StreamSpecification struct {
	StreamEnabled  bool
	StreamViewType string
}
