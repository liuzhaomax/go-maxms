package utils

import (
	"strconv"
)

func Str2Uint32(str string) (uint32, error) {
	if str == "" {
		str = "0"
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}
