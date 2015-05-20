package dynago

import "github.com/underarmour/dynago/schema"

// A Mock executor
type MockExecutor struct {
	PutItemCalled          bool
	PutItemCalledWithTable string
	PutItemCalledWithItem  Document
	PutItemResult          *PutItemResult
	PutItemError           error

	GetItemCalled          bool
	GetItemCalledWithTable string
	GetItemCalledWithKey   Document
	GetItemResultItem      Document
	GetItemError           error

	QueryCalled                              bool
	QueryCalledWithTable                     string
	QueryCalledWithIndexName                 string
	QueryCalledWithKeyConditionExpression    string
	QueryCalledWithExpressionAttributeValues Document
	QueryError                               error
	QueryResultItems                         []Document
}

func (e *MockExecutor) GetItem(getItem *GetItem) (*GetItemResult, error) {
	e.GetItemCalled = true
	e.GetItemCalledWithTable = getItem.req.TableName
	e.GetItemCalledWithKey = getItem.req.Key
	return &GetItemResult{Item: e.GetItemResultItem}, e.GetItemError
}

func (e *MockExecutor) PutItem(putItem *PutItem) (*PutItemResult, error) {
	e.PutItemCalled = true
	e.PutItemCalledWithTable = putItem.req.TableName
	e.PutItemCalledWithItem = putItem.req.Item
	return e.PutItemResult, e.PutItemError
}

func (e *MockExecutor) Query(query *Query) (*QueryResult, error) {
	e.QueryCalled = true
	e.QueryCalledWithTable = query.req.TableName
	e.QueryCalledWithIndexName = query.req.IndexName
	e.QueryCalledWithKeyConditionExpression = query.req.KeyConditionExpression
	e.QueryCalledWithExpressionAttributeValues = query.req.ExpressionAttributeValues
	return &QueryResult{Items: e.QueryResultItems}, e.QueryError
}

func (e *MockExecutor) UpdateItem(*UpdateItem) (*UpdateItemResult, error) {
	return nil, nil
}

func (e *MockExecutor) CreateTable(*schema.CreateRequest) (*schema.CreateResponse, error) {
	return nil, nil
}

func (e *MockExecutor) BatchWriteItem(*BatchWrite) (*BatchWriteResult, error) {
	return nil, nil
}
