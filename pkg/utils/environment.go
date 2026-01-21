package utils

import "strings"

func IsKoyebHost(host string) bool {
	return strings.Contains(host, "koyeb.app")
}
