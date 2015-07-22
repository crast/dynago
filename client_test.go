package dynago

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"gopkg.in/underarmour/dynago.v1/internal/aws"
)

func setUp(t *testing.T) (*assert.Assertions, *Client, *MockExecutor) {
	t.Parallel()
	executor := &MockExecutor{}
	return assert.New(t), NewClient(executor), executor
}

func TestQueryParams(t *testing.T) {
	assert, client, _ := setUp(t)
	q := client.Query("Foo").Param(":start", 9)
	assert.Equal(1, len(q.req.ExpressionAttributeValues))
	q2 := q.Params(Param{":end", 4}, Param{":other", "hello"}, Param{"#name", "Name"})
	assert.Equal(3, len(q2.req.ExpressionAttributeValues))
	assert.Equal(1, len(q2.req.ExpressionAttributeNames))
	assert.Equal("Name", q2.req.ExpressionAttributeNames["#name"])
}

func TestNewAwsClient(t *testing.T) {
	assert := assert.New(t)
	client := NewAwsClient("us-east-1", "abc", "def")
	assert.IsType(&AwsExecutor{}, client.executor)
	executor := client.executor.(*AwsExecutor)
	requester := executor.Requester.(*aws.RequestMaker)
	assert.Equal("https://dynamodb.us-east-1.amazonaws.com/", requester.Endpoint)

	signer := requester.Signer.(*aws.AwsSigner)
	assert.Equal("us-east-1", signer.Region)
	assert.Equal("abc", signer.AccessKey)
	assert.Equal("def", signer.SecretKey)
}
