package core

// config dir
const configDir = "environment/config"

// logger
const SUCCESS = "成功"
const FAILURE = "失败"

// vault
const (
	KV     = "kv"
	PWD    = "pwd"
	RSA    = "rsa"
	JWT    = "jwt"
	SECRET = "secret"
	SALT   = "salt"
	PUK    = "puk"
	PRK    = "prk"
	ID     = "id"
	APP    = "app"
)

const (
	EmptyString = ""
	Bearer      = "Bearer "
	UserID      = "userID"
	// headers params
	Method        = "method"
	URI           = "uri"
	ClientIp      = "client_ip"
	UserAgent     = "user-agent"
	TraceId       = "trace_id"
	SpanId        = "span_id"
	ParentId      = "parent_id"
	Authorization = "authorization"
	AppId         = "app_id"
	RequestId     = "request_id"
	UserId        = "user_id"
	RequestURI    = "request_uri"
	// redis
	Signature = "signature"
	// jaeger
	Tracer      = "tracer"
	Parent      = "parent"
	UberTraceId = "uber-trace-id"
)
