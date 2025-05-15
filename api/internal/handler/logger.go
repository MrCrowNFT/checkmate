package handler

import (
	"checkmate/api/internal/utils"
	"context"

	log "github.com/sirupsen/logrus"
)

// creates a structured logger with request ID
func GetLoggerWithRequestID(ctx context.Context, handlerName string) *log.Entry {
	requestID := utils.GetRequestIDFromContext(ctx)

	return log.WithFields(log.Fields{
		"func":       handlerName,
		"request_id": requestID,
	})
}
