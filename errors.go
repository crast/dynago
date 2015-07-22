package dynago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/underarmour/dynago.v1/internal/dynamodb"
)

// Encapsulates errors coming from amazon/dynamodb
type Error struct {
	Type          AmazonError    // Parsed and mapped down type
	AmazonRawType string         // Raw error type from amazon
	Exception     string         // Exception from amazon
	Message       string         // Raw message from amazon
	Request       *http.Request  // If available, HTTP request
	RequestBody   []byte         // If available, raw request body bytes
	Response      *http.Response // If available, HTTP response
	ResponseBody  []byte         // If available, raw response body bytes
}

func (e *Error) Error() string {
	exception := e.Exception
	if exception == "" {
		exception = e.AmazonRawType
	}
	return fmt.Sprintf("dynago.Error(%d): %s: %s", e.Type, exception, e.Message)
}

// Parse and create the error
func (e *Error) parse(input *inputError) {
	e.AmazonRawType = input.AmazonRawType
	e.Message = input.Message
	parts := strings.Split(e.AmazonRawType, "#")
	if len(parts) >= 2 {
		e.Exception = parts[1]
		if conf, ok := amazonErrorMap[e.Exception]; ok {
			e.Type = conf.mappedError
		}
	}
}

func buildError(req *http.Request, body []byte, response *http.Response, respBody []byte) error {
	e := &Error{
		Request:      req,
		RequestBody:  body,
		Response:     response,
		ResponseBody: respBody,
	}
	dest := &inputError{}
	if err := json.Unmarshal(respBody, dest); err == nil {
		e.parse(dest)
	} else {
		e.Message = err.Error()
	}
	return e
}

type inputError struct {
	AmazonRawType string `json:"__type"`
	Message       string `json:"message"`
}

/*
AmazonError is an enumeration of error categories that Dynago returns.

There are many more actual DynamoDB errors, but many of them are redundant from
the perspective of application logic; using these mapped-down errors is a handy
way to implement the logic you want without having a really long set of switch
statements.
*/
type AmazonError int

const (
	ErrorUnknown AmazonError = iota

	ErrorConditionFailed        // Conditional put/update failed; condition not met
	ErrorCollectionSizeExceeded // Item collection (local secondary index) too large
	ErrorThroughputExceeded     // Exceeded provisioned throughput for table or shard
	ErrorNotFound               // Resource referenced by key not found
	ErrorInternalFailure        // Internal server error
	ErrorAuth                   // Encapsulates various authorization errors
	ErrorInvalidParameter       // Encapsulates many forms of invalid input errors
	ErrorServiceUnavailable     // Amazon service unavailable
	ErrorThrottling             // Amazon is throttling us, try later
	ErrorResourceInUse          // Tried to create existing table, delete a table in CREATING state, etc.
)

type amazonErrorConfig struct {
	amazonCode     string
	expectedStatus int
	mappedError    AmazonError
}

var amazonErrors []amazonErrorConfig

var amazonErrorMap map[string]*amazonErrorConfig

func init() {
	amazonErrors = make([]amazonErrorConfig, len(dynamodb.MappedErrors))
	amazonErrorMap = make(map[string]*amazonErrorConfig, len(amazonErrors))
	for i, conf := range dynamodb.MappedErrors {
		amazonErrors[i] = amazonErrorConfig{conf.AmazonCode, conf.ExpectedStatus, AmazonError(conf.MappedError)}
		amazonErrorMap[conf.AmazonCode] = &amazonErrors[i]
	}
}
