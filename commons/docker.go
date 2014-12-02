package bazooka

import "strings"

func SafeDockerString(s string) string {
	safeString := strings.Replace(s, ":", "_", -1)
	safeString = strings.Replace(safeString, "/", "_", -1)
	return safeString
}

func ShortSHA1(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}
