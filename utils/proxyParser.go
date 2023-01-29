package utils

import (
	"sooprox/types"
	"strings"
)

func ParseProxy(s string) types.Proxy {
	parts := strings.Split(s, "::")
	return types.Proxy{
		Prefix: parts[0],
		Host:   parts[1],
	}
}
