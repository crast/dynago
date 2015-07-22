Dynago
======

Dynago is a DynamoDB client API for Go.

This attempts to be a really simple, principle of least-surprise client for the DynamoDB API.

Key design tenets of Dynago:

 * Most actions are done via chaining to build filters and conditions
 * all objects are completely safe for passing between goroutines (even queries and the like)
 * To make understanding easier via docs, we use amazon's naming wherever possible.

Installation
------------
Install using `go get`:

    go get gopkg.in/underarmour/dynago.v1

Docs are at http://godoc.org/gopkg.in/underarmour/dynago.v1

Example
-------

Run a query:

```go
client := dynago.NewClient(endpoint, accessKey, secretKey)

query := client.Query(table).
	KeyConditionExpression("UserId = :uid", dynago.Param{":uid", 42}).
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
```

Type Marshaling
---------------

Dynago attempts to let you use go types instead of having to understand a whole lot about dynamo's internal type system.

Example:

```go
doc := dynago.Document{
	"name": "Bob",
	"age": 45,
	"height": 2.1,
	"address": dynago.Document{
		"city": "Boston",
	},
	"tags": dynago.StringSet{"male", "middle_aged"},
}
client.PutItem("person", doc).Execute()
```

 * Strings use golang `string`
 * Numbers can be input as `int` or `float64` but always are returned as `dynago.Number` to not lose precision.
 * Maps can be either `map[string]interface{}` or `dynago.Document`
 * Opaque binary data can be put in `[]byte`
 * String sets, number sets, binary sets are supported using `dynago.StringSet` `dynago.NumberSet` `dynago.BinarySet`
 * Lists are supported using `dynago.List`
 * `time.Time` is only accepted if it's a UTC time, and is marshaled to a dynamo string in iso8601 compact format. It comes back as a string, an can be got back using `GetTime()` on `Document`.
