package collection_helper

// IsEmpty 是否Empty
func IsEmpty[InputType any](input []InputType) bool {
	if input == nil || len(input) == 0 {
		return true
	}
	return false
}

// IsNotEmpty 是否为Empty
func IsNotEmpty[InputType any](input []InputType) bool {
	return !IsEmpty(input)
}
