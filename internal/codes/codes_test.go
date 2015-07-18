package codes_test

import (
	"testing"

	"github.com/underarmour/dynago"
	"github.com/underarmour/dynago/internal/codes"

	"github.com/stretchr/testify/assert"
)

func TestCodeMatch(t *testing.T) {
	check := func(a dynago.AmazonError, b codes.ErrorCode) {
		assert.Equal(t, int(a), int(b))
	}
	check(dynago.ErrorUnknown, codes.ErrorUnknown)
	check(dynago.ErrorConditionFailed, codes.ErrorConditionFailed)
	check(dynago.ErrorCollectionSizeExceeded, codes.ErrorCollectionSizeExceeded)
	check(dynago.ErrorThroughputExceeded, codes.ErrorThroughputExceeded)
	check(dynago.ErrorNotFound, codes.ErrorNotFound)
	check(dynago.ErrorInternalFailure, codes.ErrorInternalFailure)
	check(dynago.ErrorAuth, codes.ErrorAuth)
	check(dynago.ErrorInvalidParameter, codes.ErrorInvalidParameter)
	check(dynago.ErrorServiceUnavailable, codes.ErrorServiceUnavailable)
	check(dynago.ErrorThroughputExceeded, codes.ErrorThroughputExceeded)
	check(dynago.ErrorNotFound, codes.ErrorNotFound)
	check(dynago.ErrorAuth, codes.ErrorAuth)
	check(dynago.ErrorInvalidParameter, codes.ErrorInvalidParameter)
	check(dynago.ErrorServiceUnavailable, codes.ErrorServiceUnavailable)
	check(dynago.ErrorThrottling, codes.ErrorThrottling)
	check(dynago.ErrorResourceInUse, codes.ErrorResourceInUse)

	check(dynago.ErrorExpiredIterator, codes.ErrorExpiredIterator)
	check(dynago.ErrorTrimmedData, codes.ErrorTrimmedData)
}
