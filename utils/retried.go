package utils

import (
	"github.com/avast/retry-go"
	"time"
)

func NewRetriedOperation(retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,

		// reduce attempts to 5 in order to not mess around with the API
		retry.Attempts(5),

		// increased delay to 2s in order to give sufficient backoff
		retry.Delay(time.Second * 2),

		// increased random to 3s in order to give sufficient room for random
		retry.MaxJitter(time.Second * 3),
	)
}
