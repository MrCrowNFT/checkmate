package utils

import (
	"context"

	"github.com/google/uuid"
)

// custom type for context keys to avoid collisions
type contextKey string

// key used to store and retrieve request IDs from context
var RequestIDKey contextKey

// initializes the request ID key
func InitRequestIDKey() {
	RequestIDKey = "request_id"
}

func GenerateRequestID() string {
	return uuid.New().String()
}

// retrieves the request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return "unknown"
}
