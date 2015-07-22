package dynago_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/underarmour/dynago.v1"
)

func TestNumberIntValReturnsTheValueAsAnInt(t *testing.T) {
	num := dynago.Number("18")
	intVal, err := num.IntVal()
	assert.Equal(t, 18, intVal)
	assert.Nil(t, err)
}

func TestNumberIntValReturnsAnErrorIfItCannotParseTheValue(t *testing.T) {
	num := dynago.Number("nope")
	intVal, err := num.IntVal()
	assert.Equal(t, 0, intVal)
	assert.Error(t, err)
}

func TestNumberInt64ValReturnsTheValueAsAnInt(t *testing.T) {
	num := dynago.Number("18")
	intVal, err := num.Int64Val()
	assert.Equal(t, int64(18), intVal)
	assert.Nil(t, err)
}

func TestNumberUint64ValReturnsTheValueAsAnInt(t *testing.T) {
	num := dynago.Number("123456789012")
	intVal, err := num.Uint64Val()
	assert.Equal(t, uint64(123456789012), intVal)
	assert.Nil(t, err)
}

func TestNumberInt64ValReturnsAnErrorIfItCannotParseTheValue(t *testing.T) {
	num := dynago.Number("nope")
	intVal, err := num.Int64Val()
	assert.Equal(t, int64(0), intVal)
	assert.Error(t, err)
}

func TestNumberFloatValReturnsTheValueAsAnfloat(t *testing.T) {
	num := dynago.Number("18.12")
	floatVal, err := num.FloatVal()
	assert.Equal(t, float64(18.12), floatVal)
	assert.Nil(t, err)
}

func TestNumberFloatValReturnsAnErrorIfItCannotParseTheValue(t *testing.T) {
	num := dynago.Number("nope")
	floatVal, err := num.FloatVal()
	assert.Equal(t, float64(0), floatVal)
	assert.Error(t, err)
}

func TestListAsDocumentListReturnsAListOfDocuments(t *testing.T) {
	list := dynago.List{dynago.Document{"id": 1}}
	docList, _ := list.AsDocumentList()
	assert.Equal(t, dynago.Document{"id": 1}, docList[0])
}

func TestListAsDocumentListReturnsAnErrorIfThereAreNonDocuments(t *testing.T) {
	list := dynago.List{dynago.Document{"real": "item"}, "imnotadocument"}
	_, err := list.AsDocumentList()
	assert.Equal(t, err.Error(), "item at index 1 was not a Document")
}

func TestListAsDocumentListReturnsTheDocumentsUpToTheFirstNonDocument(t *testing.T) {
	list := dynago.List{dynago.Document{"real": "item"}, "imnotadocument", dynago.Document{"i won't": "show up"}}
	docList, _ := list.AsDocumentList()
	assert.Equal(t, []dynago.Document{dynago.Document{"real": "item"}}, docList)
}

func TestDocumentGetStringReturnsTheUnderlyingValueAsAString(t *testing.T) {
	doc := dynago.Document{"name": "Timmy Testerson"}
	assert.Equal(t, "Timmy Testerson", doc.GetString("name"))
}

func TestDocumentGetStringReturnsAnEmptyStringWhenTheKeyIsNotPresent(t *testing.T) {
	doc := dynago.Document{}
	assert.Equal(t, "", doc.GetString("name"))
}

func TestDocumentGetNumberReturnsTheDynagoNumberWrappingTheValue(t *testing.T) {
	doc := dynago.Document{"id": dynago.Number("12")}
	assert.Equal(t, dynago.Number("12"), doc.GetNumber("id"))
}

func TestDocumentGetNumberReturnsAnEmptyNumberWhenTheKeyIsNotPresent(t *testing.T) {
	doc := dynago.Document{}
	assert.Equal(t, dynago.Number(""), doc.GetNumber("id"))
}

func TestDocumentGetNumberPanicsIfTheUnderlyingTypeIsNotANumber(t *testing.T) {
	doc := dynago.Document{"id": "not-a-dynago-number"}
	assert.Panics(t, func() {
		doc.GetNumber("id")
	})
}

func TestDocumentGetStringSetReturnsTheStringSetValue(t *testing.T) {
	doc := dynago.Document{"vals": dynago.StringSet{"val1", "val2"}}
	assert.Equal(t, dynago.StringSet{"val1", "val2"}, doc.GetStringSet("vals"))
}

func TestDocumentGetStringSetReturnsAnEmptyStringSetWhenTheKeyDoesNotExist(t *testing.T) {
	doc := dynago.Document{}
	assert.Equal(t, dynago.StringSet{}, doc.GetStringSet("vals"))
}

func TestDocumentGetStringSetPanic(t *testing.T) {
	doc := dynago.Document{"vals": "not-a-string-slice"}
	assert.Panics(t, func() {
		doc.GetStringSet("vals")
	})
}

func TestDocumentGetTimeReturnsTheTimeValueFromISO8601(t *testing.T) {
	doc := dynago.Document{"time": "1990-04-16T00:00:00Z"}
	val, _ := time.Parse("2006-01-02T15:04:05Z", "1990-04-16T00:00:00Z")
	assert.Equal(t, &val, doc.GetTime("time"))
}

func TestDocumentGetTimeReturnsNilWhenTheKeyDoesNotExist(t *testing.T) {
	doc := dynago.Document{}
	assert.Nil(t, doc.GetTime("time"))
}

func TestDocumentGetTimePanicsWhenFormatIsNotIso8601(t *testing.T) {
	doc := dynago.Document{"time": "Foo"}
	assert.Panics(t, func() { doc.GetTime("time") })
}

func TestDocumentMarshalJSONDoesNotIncludeEmptyValues(t *testing.T) {
	doc := dynago.Document{"key1": "shows up", "key2": 9, "fields": dynago.StringSet([]string{"is", "present"}), "id": "", "name": nil, "tags": []string{}}
	jsonDoc, _ := doc.MarshalJSON()

	assert.Contains(t, string(jsonDoc), `"fields":{"SS":["is","present"]}`)
	assert.Contains(t, string(jsonDoc), `"key1":{"S":"shows up"}`)
	assert.Contains(t, string(jsonDoc), `"key2":{"N":"9"}`)
}

func TestDocumentGetBool(t *testing.T) {
	doc := dynago.Document{"val": dynago.Number("1")}
	assert.Equal(t, true, doc.GetBool("val"))

	doc = dynago.Document{}
	assert.Equal(t, false, doc.GetBool("val"))

	doc = dynago.Document{"val": nil}
	assert.Equal(t, false, doc.GetBool("val"))

	doc = dynago.Document{"val": dynago.Number("b")}
	assert.Panics(t, func() {
		doc.GetBool("val")
	})

	doc = dynago.Document{"val": "hello"}
	assert.Panics(t, func() {
		doc.GetBool("val")
	})
}

func TestDocumentGetList(t *testing.T) {
	doc := dynago.Document{"vals": dynago.List{"val1", "val2"}, "wrongtype": 4}
	assert.Equal(t, dynago.List{"val1", "val2"}, doc.GetList("vals"))
	assert.Equal(t, dynago.List(nil), doc.GetList("notarealkey"))
	assert.Panics(t, func() { doc.GetList("wrongtype") })
}
