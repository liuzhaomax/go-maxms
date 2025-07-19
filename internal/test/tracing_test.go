package test_test

import (
	"fmt"
	"testing"

	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

func TestTraceID(t *testing.T) {
	fmt.Println(ext.TraceID())
	fmt.Println(ext.SpanID())
}
