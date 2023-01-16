package syncthing

import (
	"strings"
)

func ParsePortFromAddress(addr string) string {
	if strings.Contains(addr, ":") {
		lst := strings.Split(addr, ":")
		return lst[len(lst)-1]
	}
	return ""
}
