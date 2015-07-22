package dynago_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"gopkg.in/underarmour/dynago.v1"
	"gopkg.in/underarmour/dynago.v1/schema"
)

type functional struct {
	client *dynago.Client
}

func (f *functional) setUp(t *testing.T) (*assert.Assertions, *dynago.Client) {
	if testing.Short() {
		t.SkipNow()
	}
	dynago.DebugFunc = t.Logf

	if f.client == nil {
		endpoint := os.Getenv("DYNAGO_TEST_ENDPOINT")
		if endpoint == "" {
			t.SkipNow()
		}
		executor := dynago.NewAwsExecutor(endpoint, "us-east-1", "AKIAEXAMPLE", "SECRETEXAMPLE")
		f.client = dynago.NewClient(executor)
		makeTables(t, f.client)
	}
	return assert.New(t), f.client
}

var funcTest functional

func makeTables(t *testing.T, client *dynago.Client) {
	complexIndexed := complexIndexedSchema()
	hashTable := schema.NewCreateRequest("Person").HashKey("Id", schema.Number)
	hashRange := schema.NewCreateRequest("Posts").
		HashKey("UserId", schema.Number).
		RangeKey("Dated", schema.Number)

	tables := []*schema.CreateRequest{hashTable, hashRange, complexIndexed}
	for _, table := range tables {
		_, err := client.CreateTable(table)
		if err != nil {
			if e, ok := err.(*dynago.Error); ok && e.Type == dynago.ErrorResourceInUse {
				continue
			}
			panic(err)
		}
	}
}

func complexIndexedSchema() *schema.CreateRequest {
	return &schema.CreateRequest{
		TableName: "Indexed",
		AttributeDefinitions: []schema.AttributeDefinition{
			{"Id", schema.String},
			{"UserId", schema.String},
			{"Dated", schema.Number},
		},
		KeySchema: []schema.KeySchema{
			{"UserId", schema.HashKey},
			{"Id", schema.RangeKey},
		},
		ProvisionedThroughput: schema.ProvisionedThroughput{1, 1},
		GlobalSecondaryIndexes: []schema.SecondaryIndex{
			{
				IndexName:  "index1",
				Projection: schema.Projection{schema.ProjectAll, nil},
				KeySchema: []schema.KeySchema{
					{"Id", schema.HashKey},
				},
				ProvisionedThroughput: schema.ProvisionedThroughput{1, 1},
			},
		},
		LocalSecondaryIndexes: []schema.SecondaryIndex{
			{
				IndexName:  "index2",
				Projection: schema.Projection{schema.ProjectInclude, []string{"Foo", "Bar"}},
				KeySchema: []schema.KeySchema{
					{"UserId", schema.HashKey},
					{"Dated", schema.RangeKey},
				},
				ProvisionedThroughput: schema.ProvisionedThroughput{1, 1},
			},
		},
	}
}

func TestDeleteItem_functional(t *testing.T) {
	assert, client := funcTest.setUp(t)
	_, err := client.PutItem("Person", person(47, "Mary")).Execute()
	assert.NoError(err)

	key := dynago.HashKey("Id", 47)
	di := client.DeleteItem("Person", key).
		ConditionExpression("#n <> :name", dynago.Param{"#n", "Name"}, dynago.Param{":name", "Mary"}).
		ReturnValues(dynago.ReturnAllOld)
	result, err := di.Execute()
	assert.Nil(result)
	assert.NotNil(err)
	e := err.(*dynago.Error)
	assert.Equal(dynago.ErrorConditionFailed, e.Type)

	result, err = di.ConditionExpression("#n <> :name", dynago.Param{":name", "Albert"}).Execute()
	assert.NoError(err)
	assert.NotNil(result)
	doc := dynago.Document{"Name": "Mary", "IncVal": dynago.Number("1"), "Id": dynago.Number("47")}
	assert.Equal(doc, result.Attributes)

	result, err = di.ReturnValues(dynago.ReturnNone).Execute()
	assert.Nil(err)
	assert.Nil(result)
}

func TestGet(t *testing.T) {
	assert, client := funcTest.setUp(t)
	putResp, err := client.PutItem("Person", person(42, "Bob")).Execute()
	assert.NoError(err)
	assert.Nil(putResp)

	response, err := client.GetItem("Person", dynago.HashKey("Id", 42)).Execute()
	assert.Equal("Bob", response.Item["Name"])
	assert.IsType(dynago.Number("1"), response.Item["Id"])
	assert.Equal(dynago.Number("42"), response.Item["Id"])
}

func TestPutItemConditional(t *testing.T) {
	assert, client := funcTest.setUp(t)
	doc := person(45, "Joe")
	doc["Count"] = 94
	client.PutItem("Person", doc).Execute()

	doc["Count"] = 45

	basePut := client.PutItem("Person", doc).
		ConditionExpression("#c > :val").
		Param("#c", "Count")

	_, err := basePut.Param(":val", 100).Execute()

	e := err.(*dynago.Error)
	assert.Equal(dynago.ErrorConditionFailed, e.Type)

	_, err = basePut.Param(":val", 50).Execute()
	assert.NoError(err)
}

func TestBatchWrite(t *testing.T) {
	assert, client := funcTest.setUp(t)
	_, err := client.PutItem("Person", person(4, "ToDelete")).Execute()
	assert.NoError(err)

	p1 := person(1, "Joe")
	p2 := person(2, "Mary")
	p3 := person(3, "Amy")
	_, err = client.BatchWrite().
		Put("Person", p1, p2, p3).
		Delete("Person", dynago.HashKey("Id", 4)).
		Execute()

	assert.NoError(err)

	response, err := client.GetItem("Person", dynago.HashKey("Id", 2)).Execute()
	assert.Equal("Mary", response.Item["Name"])
}

func TestUpdateItemConditional(t *testing.T) {
	assert, client := funcTest.setUp(t)
	_, err := client.PutItem("Person", person(5, "ToUpdate")).Execute()
	assert.NoError(err)
	result, err := client.UpdateItem("Person", dynago.HashKey("Id", 5)).
		UpdateExpression("SET #n = :name").
		Param("#n", "Name").Param(":name", "Bob").
		ReturnValues(dynago.ReturnUpdatedNew).Execute()
	assert.NoError(err)
	assert.NotNil(result)
	assert.NotNil(result.Attributes)
	assert.Equal("Bob", result.Attributes["Name"])
}

func TestTableActions(t *testing.T) {
	tables := []string{"abc", "def", "ghi"}
	assert, client := funcTest.setUp(t)
	for _, name := range tables {
		_, err := client.CreateTable(schema.NewCreateRequest(name).HashKey("Id", schema.Number))
		if e, ok := err.(*dynago.Error); !ok || e.Type != dynago.ErrorResourceInUse {
			assert.NoError(err)
		}
	}
	list, err := client.ListTables().Limit(10).Execute()
	assert.NoError(err)
	assert.NotNil(list)
	assert.True(len(list.TableNames) > len(tables))
	assert.Nil(list.Next())

	// Pagination of tables should work
	list1, err := client.ListTables().Limit(2).Execute()
	assert.NoError(err)
	assert.NotNil(list1.Next())
	assert.Equal(2, len(list1.TableNames))
	assert.Equal(list.TableNames[:2], list1.TableNames[:2])
	list2, err := list1.Next().Execute()
	assert.Equal(list.TableNames[2:4], list2.TableNames[:2])

	for _, name := range tables {
		_, err := client.DeleteTable(name)
		assert.NoError(err)
	}

}

func TestDescribeTable(t *testing.T) {
	assert, client := funcTest.setUp(t)
	response, err := client.DescribeTable("Bogus")
	assert.Error(err)

	response, err = client.DescribeTable("Posts")
	assert.NoError(err)
	assert.Equal("Posts", response.Table.TableName)
	assert.Equal(2, len(response.Table.AttributeDefinitions))
	assert.Equal(0, len(response.Table.GlobalSecondaryIndexes))

	response, err = client.DescribeTable("Indexed")
	assert.NoError(err)
	assert.Equal(1, len(response.Table.GlobalSecondaryIndexes))
	assert.Equal(1, len(response.Table.LocalSecondaryIndexes))
	lsi := response.Table.LocalSecondaryIndexes[0]
	assert.Equal("index2", lsi.IndexName)
	assert.Equal(2, len(lsi.KeySchema))
	assert.Equal(schema.ProjectInclude, lsi.Projection.ProjectionType)
	assert.Equal([]string{"Foo", "Bar"}, lsi.Projection.NonKeyAttributes)
}

func TestPutItemReturnValues(t *testing.T) {
	assert, client := funcTest.setUp(t)
	doc := dynago.Document{"UserId": 50, "Dated": 2, "Title": "abc"}
	response, err := client.PutItem("Posts", doc).Execute()
	assert.Nil(response)
	assert.NoError(err)

	// Now test return values
	doc["Title"] = "def"
	response, err = client.PutItem("Posts", doc).
		ReturnValues(dynago.ReturnAllOld).
		Execute()
	assert.NoError(err)
	assert.NotNil(response)
	assert.Equal("abc", response.Attributes["Title"])
}

func TestQueryPagination(t *testing.T) {
	assert, client := funcTest.setUp(t)
	assert.NoError(batchAddPosts(client, 42, 100, 118))

	// Paginate the posts
	q := client.Query("Posts").
		KeyConditionExpression("UserId = :uid", dynago.Param{":uid", 42}).
		FilterExpression("Dated <> :d", dynago.Param{":d", 101}).
		Limit(10)
	results, err := q.Execute()
	assert.NoError(err)
	assert.Equal(9, len(results.Items))
	assert.Equal(dynago.Number("42"), results.Items[0]["UserId"])
	assert.Equal(dynago.Number("100"), results.Items[0]["Dated"])

	// Check that we skipped something.
	assert.Equal(10, results.ScannedCount)
	assert.Equal(9, results.Count)

	assert.NotNil(results.LastEvaluatedKey)
	assert.Equal(2, len(results.LastEvaluatedKey))
	assert.NotNil(results.Next())

	// page 2, also use ProjectionExpression
	results, err = results.Next().ProjectionExpression("Dated").Execute()
	assert.NoError(err)
	assert.Equal(8, len(results.Items))
	assert.Equal(1, len(results.Items[0]))
	assert.Nil(results.LastEvaluatedKey)
	assert.Nil(results.Next())
}

func TestScanBasic(t *testing.T) {
	assert, client := funcTest.setUp(t)
	assert.NoError(batchAddPosts(client, 50, 100, 120))
	scan := client.Scan("Posts").Limit(20)
	result, err := scan.Execute()
	assert.NoError(err)
	assert.NotEqual(0, len(result.Items))
	assert.NotNil(result.LastEvaluatedKey)
	assert.NotNil(result.Next())
	result2, err := scan.ExclusiveStartKey(result.LastEvaluatedKey).Execute()
	assert.NotEqual(result.Items, result2.Items)
	assert.Nil(result2.Next())

	// ensure Next method matches manual ExclusiveStartKey
	result2a, err := result.Next().Execute()
	assert.NoError(err)
	assert.Equal(result2a.Items, result2.Items)
}

func person(id int, name string) dynago.Document {
	return dynago.Document{"Id": id, "Name": name, "IncVal": 1}
}

func post(uid int, dated int) dynago.Document {
	return dynago.Document{"UserId": uid, "Dated": dated}
}

func batchAddPosts(client *dynago.Client, userId, start, end int) error {
	// Add some posts
	writer := client.BatchWrite()
	for i := start; i < end; i++ {
		writer = writer.Put("Posts", dynago.Document{"UserId": userId, "Dated": i})
	}
	_, err := writer.Execute()
	return err
}
