package dynago_test

import (
	"fmt"
	"gopkg.in/underarmour/dynago.v1"
)

var region, accessKey, secretKey, table string

func Example() {
	client := dynago.NewAwsClient(region, accessKey, secretKey)

	query := client.Query(table).
		FilterExpression("NumViews > :views").
		Param(":views", 50).
		Desc()

	result, err := query.Execute()
	if err != nil {
		// do something
	}
	for _, row := range result.Items {
		fmt.Printf("Name: %s, Views: %d", row["Name"], row["NumViews"])
	}
}

func Example_marshaling(client *dynago.Client) {
	type MyStruct struct {
		Id          int64
		Name        string
		Description string
		Tags        []string
		Address     struct {
			City  string
			State string
		}
	}

	var data MyStruct

	doc := dynago.Document{
		// Basic fields like numbers and strings get marshaled automatically
		"Id":          data.Id,
		"Name":        data.Name,
		"Description": data.Description,
		// StringSet is compatible with []string so we can simply cast it
		"Tags": dynago.StringSet(data.Tags),
		// We don't automatically marshal structs, nest it in a document
		"Address": dynago.Document{
			"City":  data.Address.City,
			"State": data.Address.State,
		},
	}

	client.PutItem("Table", doc).Execute()
}
