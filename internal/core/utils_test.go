package core

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
)

func TestGetFuncName(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "测试返回函数名称功能",
			want: "core.TestGetFuncName",
		},
	}
	funcName := GetFuncName()
	gotSlice := strings.Split(funcName, "/")
	got := gotSlice[len(gotSlice)-1]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, got, tc.want, fmt.Sprintf("\n*** Expected: \n %#v \n*** Got: \n %#v", tc.want, got))
		})
	}
}
