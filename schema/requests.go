package schema

type CreateRequest struct {
	TableName             string
	AttributeDefinitions  []AttributeDefinition
	KeySchema             []KeySchema
	ProvisionedThroughput ProvisionedThroughput
	// TODO local and global secondary indexes
	GlobalSecondaryIndexes []SecondaryIndex
}

type SecondaryIndex struct {
	IndexName             string
	KeySchema             []KeySchema
	Projection            Projection
	ProvisionedThroughput ProvisionedThroughput
}

type Projection struct {
	ProjectionType string
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

type CreateResponse struct {
	TableDescription TableDescription
}
