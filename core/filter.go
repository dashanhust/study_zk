package core

import "strings"

// 功能：返回没有以符号 - 开头的参数
func excludeOptions(args []string) []string {
	ret := make([]string, 0, len(args))
	for i := range args {
		if !strings.HasPrefix(args[i], "-") {
			ret = append(ret, args[i])
		}
	}
	return ret
}
