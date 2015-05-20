package dynago

import (
	"github.com/underarmour/dynago/schema"
)

/*
Create a new dynamo client.

region is the AWS region, e.g. us-east-1.
accessKey is your amazon access key ID.
secretKey is your amazon secret key ID.
*/
func NewClient(region, accessKey, secretKey string) *Client {
	endpoint := "https://dynamodb." + region + ".amazonaws.com/"
	return NewEndpointClient(endpoint, region, accessKey, secretKey)
}

func NewEndpointClient(endpoint, region, accessKey, secretKey string) *Client {
	return &Client{
		executor: NewAwsExecutor(endpoint, region, accessKey, secretKey),
	}
}

/*
Create a new client with a custom executor.

This is mainly used for unit test and mock scenarios.
*/
func NewClientExecutor(executor Executor) *Client {
	return &Client{executor}
}

type Client struct {
	executor Executor
}

/*
Compose a batch write.

Batch writes can compose a number of put or delete, even across multiple tables, in a single operation.
*/
func (c *Client) BatchWrite() *BatchWrite {
	return newBatchWrite(c)
}

/*
Compose a GetItem on a dynamo table.

key should be a Document containing enough attributes to describe the primary key.

You can use the HashKey or HashRangeKey helpers to help build a key:

	client.GetItem("foo", dynago.HashKey("Id", 45))

	client.GetItem("foo", dynago.HashRangeKey("UserId", 45, "Date", "20140512"))
*/
func (c *Client) GetItem(table string, key Document) *GetItem {
	return newGetItem(c, table, key)
}

/*
Compose a Query on a dynamo table.

This returns a new Query struct which you can compose via chaining to build the query you want.
Then finish the chain by calling Execute() to run the query.
*/
func (c *Client) Query(table string) *Query {
	return newQuery(c, table)
}

/*
Compose a PutItem on a dynamo table.

item should be a document representing the record and containing the attributes for the primary key.

Like all the other requests, you must call `Execute()` to run this.
*/
func (c *Client) PutItem(table string, item Document) *PutItem {
	return newPutItem(c, table, item)
}

/*
Compose an UpdateItem on a dynamo table.
*/
func (c *Client) UpdateItem(table string, key Document) *UpdateItem {
	return newUpdateItem(c, table, key)
}

/*
Create a table.
*/
func (c *Client) CreateTable(req *schema.CreateRequest) (*schema.CreateResponse, error) {
	return c.executor.CreateTable(req)
}
