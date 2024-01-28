package middleware

// Headers
const (
	contentTypeHeader            = "Content-Type"
	applicationJSONValue         = "application/json"
	xForwardedForHeader          = "x-forwarded-for"
	xServiceAuthenticationHeader = "x-service-authentication"
)

// Errors
const (
	InternalServerError             string = "An unexpected error occurred while processing your request."
	RateLimitExceededError          string = "Slow down! Too many requests. Try again shortly. Thank you!"
	MissingXForwardedForHeaderError string = "X-Forwarded-For header is required for accurate processing."
)

type ReqContextKey string

const (
	ClientIP              ReqContextKey = "client_ip"
	ServiceAuthentication ReqContextKey = "service_authentication"
	AuthenticatedService  ReqContextKey = "authenticated_service"
)
