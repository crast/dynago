/*
Dynago is a DynamoDB client API for Go.

Dynago differs from other Dynamo clients for Go in that it tries to mirror
DynamoDB's core API closely: Most methods, attributes, and names are made
following Dynamo's own naming conventions. This allows it to be clear which
API is being accessed and allows finding complementary docs on Amazon's side
easier.

Filter Chaining

A key design concept is the use of chaining to build filters and conditions,
similar to some ORM frameworks. This allows using sub-features like conditional
puts, expression post-filtering, and so on to be clearer, because this means
a conditional put is simply a PutItem with a condition expression tacked on.

	query := client.Query("Table").
		KeyConditionExpression("Foo = :foo", dynago.Param{":foo", 42}).
		Limit(40).Desc()
	result, err := query.Execute()

All the various item-based query actions are evaluated when you call the
Execute() method on a filter chain.

Type Marshaling

Dynago tries to marshal to/from Go types where possible:

 * Strings use golang string
 * Numbers can be input as int, int64, float64, etc but always are returned as
   dynago.Number to not lose precision.
 * Maps can be either map[string]interface{} or dynago.Document
 * Opaque binary data can be put in []byte
 * String sets, number sets, binary sets are supported using dynago.StringSet,
   dynago.NumberSet, dynago.BinarySet
 * Lists are supported using dynago.List
 * time.Time is only accepted if it's a UTC time, and is marshaled to a dynamo
   string in iso8601 compact format. It comes back as a string, an can be
   unmarshaled back using GetTime on document.

Query Parameters

Nearly all the operations on items allow using DynamoDB's expression language to
do things like query filtering, attribute projection, and so on. In order to provide
literal values,  queries are parametric, just like many SQL engines:

	SET Foo = Foo + :incr
	DELETE Person.#n

DynamoDB has two fields it uses for parameters: ExpressionAttributeNames for name
aliases, and ExpressionAttributeValues for parametric values.  For simplicity, in the
Dynago library both of those are serviced by Param. This is okay because parameters and
aliases are non-ambiguous in that the former are named e.g. ":foo" and the latter "#foo".

So a conditional PutItem might look like:

	client.PutItem(table, item).
		ConditionExpression("Foo.#n = :fooName").
		Param("#n", "Name").Param(":fooName", "Bob").
		Execute()

In this case, we only execute the query if the value at document path Foo.Name was the
string value "Bob". Note we used the "Param" helper for setting both values.

There are also helpers that let you set multiple parameters, or inline with expressions:

	query.FilterExpression("#foo > :foo", Param{":foo", 45}, Param{"#foo", "Foo"})
	-or-
	query.Params(Param{":foo", 45}, Param{"#foo", "Foo"})
	-or-
	query.Params(Document{":foo":45, "#foo": "Foo"})

*/
package dynago // import gopkg.in/underarmour/dynago.v1
