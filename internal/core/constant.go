package core

// config dir
const configDir = "environment/config"

// logger
const SUCCESS = "成功"
const FAILURE = "失败"

// vault
const (
	Kv     = "kv"
	Pwd    = "pwd"
	Rsa    = "rsa"
	Jwt    = "jwt"
	Secret = "secret"
	Salt   = "salt"
	Puk    = "puk"
	Prk    = "prk"
)

const (
	EmptyString = ""
	// headers params
	ClientIp      = "Client_ip"
	UserAgent     = "User-Agent"
	TraceId       = "Trace_id"
	SpanId        = "Span_id"
	ParentId      = "Parent_id"
	Authorization = "Authorization"
)
