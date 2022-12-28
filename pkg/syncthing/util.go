package syncthing

import (
	"fmt"
	"strings"
)

func parsePortFromAddress(addr string) (string, error) {
	if strings.Contains(addr, ":") {
		return strings.Split(addr, ":")[1], nil
	}
	return "", fmt.Errorf("failed to parse port from address: %s", addr)
}
