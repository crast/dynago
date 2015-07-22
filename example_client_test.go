package dynago_test

import (
	"fmt"
	"gopkg.in/underarmour/dynago.v1"
	"gopkg.in/underarmour/dynago.v1/schema"
)

func ExampleClient_BatchWrite(client *dynago.Client) {
	record1 := dynago.Document{"Id": 1, "Name": "Person1"}
	record2 := dynago.Document{"Id": 2, "Name": "Person2"}

	// We can put and delete at the same time to multiple tables.
	client.BatchWrite().
		Put("Table1", record1, record2).
		Put("Table2", dynago.Document{"Name": "Other"}).
		Delete("Table2", dynago.HashKey("Id", 42)).
		Execute()
}

func ExampleClient_CreateTable_basic(client *dynago.Client) {
	// NewCreateRequest creates a table with simple defaults.
	// You can use chaining to set the hash and range keys.
	table1 := schema.NewCreateRequest("TableName").
		HashKey("UserId", schema.Number).
		RangeKey("Date", schema.String)

	table1.ProvisionedThroughput.ReadCapacityUnits = 45
	client.CreateTable(table1)
}

func ExampleClient_CreateTable_full(client *dynago.Client) {
	// Most of the time we don't need the full syntax for making create requests
	// It's shown here mostly for purpose of documentation
	req := &schema.CreateRequest{
		TableName: "PersonalPages",
		AttributeDefinitions: []schema.AttributeDefinition{
			{"UserId", schema.Number},
			{"Title", schema.String},
		},
		KeySchema: []schema.KeySchema{
			{"UserId", schema.HashKey},
			{"Title", schema.RangeKey},
		},
		ProvisionedThroughput: schema.ProvisionedThroughput{
			ReadCapacityUnits:  45,
			WriteCapacityUnits: 72,
		},
	}
	if response, err := client.CreateTable(req); err == nil {
		fmt.Printf("table created, status %s", response.TableDescription.TableStatus)
	}
}

func ExampleClient_ListTables_paging(client *dynago.Client) {
	cursor := client.ListTables().Limit(100)
	for cursor != nil {
		response, err := cursor.Execute()
		if err != nil {
			break
		}
		fmt.Printf("%v", response.TableNames)
		cursor = response.Next()
	}
}

func ExampleClient_PutItem(client *dynago.Client) {
	doc := dynago.Document{
		"Id":   42,
		"Name": "Bob",
		"Address": dynago.Document{
			"City":  "Boston",
			"State": "MA",
		},
	}
	_, err := client.PutItem("Person", doc).Execute()
	if err != nil {
		fmt.Printf("PUT failed: %v", err)
	}
}

func ExampleClient_Query(client *dynago.Client) {
	result, err := client.Query("table").
		KeyConditionExpression("Foo = :val AND begins_with(Title, :prefix)").
		Param(":val", 45).Param(":prefix", "The adventures of").
		Execute()

	if err == nil {
		for _, row := range result.Items {
			fmt.Printf("ID: %s, Foo: %d", row["Id"], row["Foo"])
		}
	}
}

func ExampleClient_Query_pagination(client *dynago.Client) {
	query := client.Query("table").
		KeyConditionExpression("Foo = :val").
		Limit(50)

	// Keep getting results in a loop until there are no more.
	for query != nil {
		result, err := query.Execute()
		if err != nil {
			break
		}
		for _, item := range result.Items {
			fmt.Printf("Result ID %d\n", item["Id"])
		}
		query = result.Next()
	}
}

func ExampleClient_Scan_parallel(client *dynago.Client) {
	numSegments := 10
	baseScan := client.Scan("Table").Limit(1000)

	// spin up numSegments goroutines each working on a table segment
	for i := 0; i < numSegments; i++ {
		go func(scan *dynago.Scan) {
			for scan != nil {
				result, _ := scan.Execute()
				// do something with result.Items
				scan = result.Next()
			}
		}(baseScan.Segment(i, numSegments))
	}
}

func ExampleClient_UpdateItem(client *dynago.Client) {
	_, err := client.UpdateItem("Person", dynago.HashKey("Id", 42)).
		UpdateExpression("SET Name=:name").
		Param(":name", "Joe").
		Execute()

	if err != nil {
		fmt.Printf("UpdateItem failed: %v", err)
	}
}

func ExampleClient_UpdateItem_atomicIncrement(client *dynago.Client, key dynago.Document) {
	result, err := client.UpdateItem("Products", key).
		UpdateExpression("SET ViewCount = ViewCount + :incr").
		Param(":incr", 5).
		ReturnValues(dynago.ReturnUpdatedNew).
		Execute()

	if err == nil {
		fmt.Printf("View count is now %d", result.Attributes["ViewCount"])
	}
}
