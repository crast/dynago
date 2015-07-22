package dynago

import (
	"gopkg.in/underarmour/dynago.v1/schema"
)

type awsSchemaExecutor struct {
	*AwsExecutor
}

func (e awsSchemaExecutor) CreateTable(req *schema.CreateRequest) (resp *schema.CreateResult, err error) {
	err = e.MakeRequestUnmarshal("CreateTable", req, &resp)
	return
}

func (e awsSchemaExecutor) DeleteTable(req *schema.DeleteRequest) (resp *schema.DeleteResult, err error) {
	err = e.MakeRequestUnmarshal("DeleteTable", req, &resp)
	return
}

func (e awsSchemaExecutor) DescribeTable(req *schema.DescribeRequest) (resp *schema.DescribeResponse, err error) {
	err = e.MakeRequestUnmarshal("DescribeTable", req, &resp)
	return
}

type ListTables struct {
	client *Client
	req    schema.ListRequest
}

func (l ListTables) Limit(limit uint) *ListTables {
	l.req.Limit = limit
	return &l
}

func (l *ListTables) Execute() (result *ListTablesResult, err error) {
	resp, err := l.client.schemaExecutor.ListTables(l)
	if err == nil {
		result = &ListTablesResult{resp.TableNames, resp.LastEvaluatedTableName, l}
	}
	return
}

func (e awsSchemaExecutor) ListTables(list *ListTables) (resp *schema.ListResponse, err error) {
	err = e.MakeRequestUnmarshal("ListTables", list.req, &resp)
	return resp, err
}

type ListTablesResult struct {
	TableNames []string
	cursor     *string
	req        *ListTables
}

// Helper to get the ListTables for the next page of listings.
// If there is not a next page, returns nil
func (r ListTablesResult) Next() *ListTables {
	if r.cursor == nil {
		return nil
	}
	return &ListTables{r.req.client, schema.ListRequest{
		ExclusiveStartTableName: *r.cursor,
		Limit: r.req.req.Limit,
	}}
}
