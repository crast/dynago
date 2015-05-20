package dynago_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/underarmour/dynago"
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

func TestDocumentGetStringSliceReturnsTheStringSetValue(t *testing.T) {
	doc := dynago.Document{"vals": dynago.StringSet{"val1", "val2"}}
	assert.Equal(t, dynago.StringSet{"val1", "val2"}, doc.GetStringSlice("vals"))
}

func TestDocumentGetStringSliceReturnsAnEmptyStringSetWhenTheKeyDoesNotExist(t *testing.T) {
	doc := dynago.Document{}
	assert.Equal(t, dynago.StringSet{}, doc.GetStringSlice("vals"))
}

func TestDocumentGetStringSlicePanic(t *testing.T) {
	doc := dynago.Document{"vals": "not-a-string-slice"}
	assert.Panics(t, func() {
		doc.GetStringSlice("vals")
	})
}

func TestDocumentGetTimeReturnsTheTimeValueFromISO8601(t *testing.T) {
	doc := dynago.Document{"time": "1990-04-16T00:00:00Z"}
	val, _ := time.Parse("2006-01-02T15:04:05Z", "1990-04-16T00:00:00Z")
	assert.Equal(t, dynago.Time(&val), doc.GetTime("time"))
}

func TestDocumentGetTimeReturnsNilWhenTheKeyDoesNotExist(t *testing.T) {
	doc := dynago.Document{}
	assert.Nil(t, doc.GetTime("time"))
}

func TestDocumentMarshalJSONDoesNotIncludeEmptyValues(t *testing.T) {
	doc := dynago.Document{"key1": "shows up", "key2": 9, "fields": []string{"is", "present"}, "id": "", "name": nil, "tags": []string{}}
	jsonDoc, _ := doc.MarshalJSON()

	assert.Contains(t, string(jsonDoc), `"fields":{"SS":["is","present"]}`)
	assert.Contains(t, string(jsonDoc), `"key1":{"S":"shows up"}`)
	assert.Contains(t, string(jsonDoc), `"key2":{"N":"9"}`)
}
