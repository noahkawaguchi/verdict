package utils

// Ref takes a value of any type and returns a pointer to it. Similar to `&` or `aws.String()`.
func Ref[T any](something T) *T { return &something }
