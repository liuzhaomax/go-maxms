package ext

import (
	"strings"

	"github.com/lithammer/shortuuid"
	uuid "github.com/satori/go.uuid"
)

func TraceID() string {
	return UUIDInUpper()
}

func SpanID() string {
	return UUIDInLower()
}

func UUIDInUpper() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}

func UUIDInLower() string {
	return strings.ToLower(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}

func ShortUUID() string {
	return shortuuid.New()
}
