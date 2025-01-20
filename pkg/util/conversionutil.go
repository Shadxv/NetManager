package util

func IntToPtr(i int) *int32 {
	i32 := int32(i)
	return &i32
}
