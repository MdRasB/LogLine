package middleware

type contextKey string

var (
	UserIDKey    = contextKey("user_id")
	SessionIDKey = contextKey("session_id")
)

const (
	RequestIDKey contextKey = "request_id"
)