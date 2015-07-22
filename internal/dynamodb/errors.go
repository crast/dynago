// package dynamodb maps information from dynamo itself, such as error tables and formats.
// This is an internal package to allow fast iteration on the abstraction
// without having to commit to the interface.
package dynamodb

import (
	"gopkg.in/underarmour/dynago.v1/internal/codes"
)

type DynamoErrorConfig struct {
	AmazonCode     string
	ExpectedStatus int
	MappedError    codes.ErrorCode
}

// This variable is mostly exposed so that we can document how errors are mapped
var MappedErrors = []DynamoErrorConfig{
	{"ConditionalCheckFailedException", 400, codes.ErrorConditionFailed},
	{"ResourceNotFoundException", 400, codes.ErrorNotFound},
	{"InternalFailure", 500, codes.ErrorInternalFailure},
	{"InternalServerError", 500, codes.ErrorInternalFailure},
	{"IncompleteSignature", 400, codes.ErrorAuth},
	{"InvalidParameterCombination", 400, codes.ErrorInvalidParameter},
	{"InvalidParameterValue", 400, codes.ErrorInvalidParameter},
	{"InvalidQueryParameter", 400, codes.ErrorInvalidParameter},
	{"InvalidSignatureException", 400, codes.ErrorAuth},
	{"ItemCollectionSizeLimitExceededException", 400, codes.ErrorCollectionSizeExceeded},
	{"MalformedQueryString", 404, codes.ErrorInvalidParameter},
	{"MissingAction", 400, codes.ErrorInvalidParameter},
	{"MissingAuthenticationToken", 403, codes.ErrorAuth},
	{"MissingParameter", 400, codes.ErrorInvalidParameter},
	{"OptInRequired", 403, codes.ErrorAuth},
	{"ProvisionedThroughputExceededException", 400, codes.ErrorThroughputExceeded},
	{"RequestExpired", 400, codes.ErrorAuth},
	{"ResourceInUseException", 400, codes.ErrorResourceInUse},
	{"ServiceUnavailable", 503, codes.ErrorServiceUnavailable},
	{"ThrottlingException", 400, codes.ErrorThrottling},
	{"ValidationError", 400, codes.ErrorInvalidParameter},
	{"ValidationException", 400, codes.ErrorInvalidParameter},
}
