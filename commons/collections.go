package bazooka

func CopyMap(source map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range source {
		dst[k] = v
	}
	return dst
}
