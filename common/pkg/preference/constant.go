package preference

type CtxKey string

const (
	// Respnose Status
	STATUS_SUCCESS string = "success"
	STATUS_ERROR   string = "error"

	// Database Type
	MYSQL    string = `mysql`
	POSTGRES string = `postgres`
	MARIADB  string = `mariadb`

	// Redis Type
	REDIS_APPS    string = "APPS"
	REDIS_LIMITER string = "LIMITER"
	REDIS_AUTH    string = "AUTH"

	// Logging Context Keys
	CONTEXT_KEY_TRACE_ID   CtxKey = "trace_id"
	CONTEXT_KEY_SPAN_ID    CtxKey = "span_id"
	CONTEXT_KEY_REQ_ID     CtxKey = "req_id"
	CONTEXT_KEY_ADDR       CtxKey = "addr"
	CONTEXT_KEY_USER_AGENT CtxKey = "user_agent"

	TRACE_ID string = "trace_id"
	SPAN_ID  string = "span_id"
	REQ_ID   string = "req_id"

	EVENT      string = "event"
	METHOD     string = "method"
	URL        string = "url"
	ADDR       string = "addr"
	STATUS     string = "status_code"
	LATENCY    string = "latency"
	USER_AGENT string = "user_agent"

	// Lang Header
	LANG_EN string = `en`
	LANG_ID string = `id`

	// Custom HTTP Header
	APP_LANG   string = `x-app-lang`
	REQUEST_ID string = `x-request-id`

	// Cache Control Header
	CacheControl        string = `cache-control`
	CacheMustRevalidate string = `must-revalidate`

	// Limiter Error Message
	FormatError  string = "Please check the format with your input."
	CommandError string = "The command of first number should > 0"
)
