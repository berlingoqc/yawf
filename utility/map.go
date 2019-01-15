package utility

// ConcatMap add the second map into the first one
func ConcatMap(m map[string]interface{}, m1 map[string]interface{}) {
	for k, v := range m1 {
		m[k] = v
	}
}
