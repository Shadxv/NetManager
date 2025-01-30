package util

func IntToPtr(i int) *int32 {
	i32 := int32(i)
	return &i32
}

func BoolToMongoStr(arg bool) string {
	if arg {
		return "enabled"
	}
	return "disabled"
}
