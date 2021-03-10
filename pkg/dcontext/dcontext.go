package dcontext

import (
	"context"
)

type key string

const (
	userIDKey    key = "userID"
	requestIDKey key = "requestID"
)

// SetUserID ContextへユーザIDを保存する
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext ContextからユーザIDを取得する
func GetUserIDFromContext(ctx context.Context) string {
	var userID string
	if ctx.Value(userIDKey) != nil {
		userID = ctx.Value(userIDKey).(string)
	}
	return userID
}

// SetRequestID ContextへリクエストIDを保存する
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestIDFromContext ContextからリクエストIDを取得する
func GetRequestIDFromContext(ctx context.Context) string {
	var requestID string
	if ctx.Value(requestIDKey) != nil {
		requestID = ctx.Value(requestIDKey).(string)
	}
	return requestID
}
