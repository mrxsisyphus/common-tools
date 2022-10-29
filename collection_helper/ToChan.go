package collection_helper

// ToChannel  将slice换为等量的channel 返回channel
// - InputType any 初始的数组元素类型
// 返回的是chan
func ToChannel[InputType comparable](input *[]InputType) chan InputType {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	res := make(chan InputType, len(raw))
	for _, item := range raw {
		res <- item
	}
	return res
}

// ToReadChannel  将slice换为等量的有缓存channel 返回channel 改channel 只读
// - InputType any 初始的数组元素类型
// 返回的是 <-chan
func ToReadChannel[InputType comparable](input *[]InputType) <-chan InputType {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	res := make(chan InputType, len(raw))
	for _, item := range raw {
		res <- item
	}
	return res
}
