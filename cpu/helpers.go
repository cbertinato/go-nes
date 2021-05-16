package cpu

import (
	"strconv"
	"strings"
)

func FormatInt16Bin(n uint16) string {
	res := strconv.FormatInt(int64(n), 2)
	if len(res) < 16 {
		var zeros []string
		for i := 0; i < 16-len(res); i++ {
			zeros = append(zeros, "0")
		}
		res = strings.Join(zeros, "") + res
	}
	return res
}

func FormatInt8Bin(n uint8) string {
	res := strconv.FormatInt(int64(n), 2)
	if len(res) < 16 {
		var zeros []string
		for i := 0; i < 16-len(res); i++ {
			zeros = append(zeros, "0")
		}
		res = strings.Join(zeros, "") + res
	}
	return res
}
