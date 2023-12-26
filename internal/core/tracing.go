package core

import (
	"github.com/satori/go.uuid"
	"strings"
)

func TraceID() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}

func SpanID() string {
	return strings.ToLower(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}
