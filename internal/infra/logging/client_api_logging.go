package logging

import (
	"log"
)

func Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64) {
	// Stub implementation for ClientApiLogging.log
	log.Printf("mock client api log: %v", request)
}
