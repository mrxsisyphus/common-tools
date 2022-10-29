package collection_helper

// Iter 普通迭代器 没办法自己暂停 本质 利用了无缓存的chain
func Iter[InputType any](input *[]InputType) <-chan InputType {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	ch := make(chan InputType)
	go func() {
		for _, inputType := range raw {
			ch <- inputType
		}
		close(ch)
	}()
	return ch
}

// GetIterator 获得迭代器 可以自己来暂停 本质 利用了无缓存的chain
func GetIterator[InputType any](input *[]InputType) *Iterator[InputType] {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	iterator, ch, stopCh := NewIterator[InputType]()

	go func() {
	L:
		for _, inputType := range raw {
			select {
			case <-stopCh:
				break L
			case ch <- inputType:
			}
		}
		close(ch)
	}()

	return iterator
}
