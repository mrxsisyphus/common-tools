package map_helper

// Iter 普通迭代器 没办法自己暂停 本质 利用了无缓存的chain
func Iter[K comparable, V any](m map[K]V) <-chan *Pair[K, V] {
	if len(m) <= 0 {
		return nil
	}
	ch := make(chan *Pair[K, V])
	go func() {
		for k, v := range m {
			ch <- NewPair[K, V](k, v)
		}
		close(ch)
	}()
	return ch
}

// GetIterator 获得迭代器 可以自己来暂停 本质 利用了无缓存的chain
func GetIterator[K comparable, V any](m map[K]V) *Iterator[*Pair[K, V]] {
	if len(m) <= 0 {
		return nil
	}
	iterator, ch, stopCh := NewIterator[*Pair[K, V]]()

	go func() {
	L:
		for k, v := range m {
			select {
			case <-stopCh:
				break L
			case ch <- NewPair[K, V](k, v):
			}
		}
		close(ch)
	}()

	return iterator
}
