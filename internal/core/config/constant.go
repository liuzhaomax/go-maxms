package config

// logger
const (
	SUCCESS = "成功"
	FAILURE = "失败"
)

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
	Bearer = "Bearer "
	UserID = "userID"
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
	Signature     = "signature"
	// jaeger
	UberTraceId = "uber-trace-id"
	// redis
	Nonce = "nonce"
	// pool
	MyWsConn = "MyWsConn"
)
