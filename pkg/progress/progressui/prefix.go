package progressui

import "strings"

func prefix(name, defaultPrefix string) string {
	if strings.Contains(name, "docker") {
		return defaultPrefix + "🐋 "
	}
	if strings.Contains(name, "local") {
		return defaultPrefix + "📦 "
	}
	if strings.Contains(name, "remote") {
		return defaultPrefix + "📡 "
	}
	return defaultPrefix
}
