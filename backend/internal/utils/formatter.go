package utils

func GetFormValue(values map[string][]string, key string) string {
	if vals, ok := values[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}
func Bool(b bool) *bool {
	return &b
}
func GetBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
func GetInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
func GetString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func GetInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}
func GetInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}
func GetBytes(b *[]byte) []byte {
	if b != nil {
		return *b
	}
	return []byte{}
}
