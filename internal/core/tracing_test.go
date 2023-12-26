package core

import (
	"fmt"
	"testing"
)

func TestTraceID(t *testing.T) {
	fmt.Println(TraceID())
	fmt.Println(SpanID())
}
