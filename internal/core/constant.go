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
	Method        = "Method"
	URI           = "URI"
	ClientIp      = "Client_ip"
	UserAgent     = "User-Agent"
	TraceId       = "Trace_id"
	SpanId        = "Span_id"
	ParentId      = "Parent_id"
	Authorization = "Authorization"
	AppId         = "App_id"
	RequestId     = "Request_id"
	UserId        = "User_id"
	RequestURI    = "Request_uri"
	// redis
	Signature = "signature"
	// jaeger
	Tracer = "tracer"
	Parent = "parent"
)
