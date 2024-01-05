package core

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
)

func TestIn(t *testing.T) {
	type input struct {
		slice   []int
		element int
	}
	cases := []struct {
		name string
		have input
		want bool
	}{
		{
			name: "字符串存在于切片",
			have: input{
				slice:   []int{1, 2, 3},
				element: 2,
			},
			want: true,
		},
		{
			name: "字符串不存在于切片",
			have: input{
				slice:   []int{1, 2, 3},
				element: 4,
			},
			want: false,
		},
	}
	for _, tc := range cases {
		got := In(tc.have.slice, tc.have.element)
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, got, tc.want, fmt.Sprintf("\n*** Expected: \n %#v \n*** Got: \n %#v", tc.want, got))
		})
	}
}

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

func TestGetCallerFileAndLine(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "测试返回函数位置",
			want: fmt.Sprintf("\033[1;34m%s\033[0m\n", "D:/workspace/Github/go-maxms/internal/core/utils_test.go:75"),
		},
	}
	got := GetCallerFileAndLine(2)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, got, tc.want, fmt.Sprintf("\n*** Expected: \n %#v \n*** Got: \n %#v", tc.want, got))
		})
	}
}
