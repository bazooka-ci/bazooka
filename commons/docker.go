package bazooka

import "strings"

func SafeDockerString(s string) string {
	safeString := strings.Replace(s, ":", "_", -1)
	safeString = strings.Replace(safeString, "/", "_", -1)
	return safeString
}
