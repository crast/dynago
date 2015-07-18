package schema

type CreateRequest struct {
	TableName              string
	AttributeDefinitions   []AttributeDefinition
	KeySchema              []KeySchema
	ProvisionedThroughput  ProvisionedThroughput
	GlobalSecondaryIndexes []SecondaryIndex
	LocalSecondaryIndexes  []SecondaryIndex
	StreamSpecification    *StreamSpecification `json:",omitempty"`
}

func NewCreateRequest(table string) *CreateRequest {
	return &CreateRequest{
		TableName:             table,
		ProvisionedThroughput: ProvisionedThroughput{1, 1},
	}
}

// Simple way to add a hash key and attribute definition in one go.
func (r *CreateRequest) HashKey(name string, attributeType AttributeType) *CreateRequest {
	r.ensureAttribute(name, attributeType)
	r.KeySchema = append(r.KeySchema, KeySchema{name, HashKey})
	return r
}

// Simple way to add a hash key and attribute definition in one go.
func (r *CreateRequest) RangeKey(name string, attributeType AttributeType) *CreateRequest {
	r.ensureAttribute(name, attributeType)
	r.KeySchema = append(r.KeySchema, KeySchema{name, RangeKey})
	return r
}

func (r *CreateRequest) ensureAttribute(name string, attributeType AttributeType) {
	for _, a := range r.AttributeDefinitions {
		if a.AttributeName == name {
			return
		}
	}
	r.AttributeDefinitions = append(r.AttributeDefinitions, AttributeDefinition{name, attributeType})
}

type CreateResult struct {
	TableDescription TableDescription
}

type DeleteRequest struct {
	TableName string
}

type DeleteResult struct {
	TableDescription TableDescription
}

type ListRequest struct {
	ExclusiveStartTableName string `json:",omitempty"`
	Limit                   uint   `json:",omitempty"`
}

type ListResponse struct {
	LastEvaluatedTableName *string
	TableNames             []string
}

type DescribeRequest struct {
	TableName string
}

type DescribeResponse struct {
	Table TableDescription
}
