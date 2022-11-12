package commonMapOperation

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/monitor1379/yagods/sets/hashset"
	"github.com/monitor1379/yagods/sets/linkedhashset"
	"github.com/monitor1379/yagods/sets/treeset"
	"github.com/monitor1379/yagods/utils"
	"github.com/mrxtryagain/common-tools/convert_helper"
	"github.com/mrxtryagain/common-tools/map_helper"
)

type MapOperation[K comparable, V any] struct {
	input map[K]V
}

func (m *MapOperation[K, V]) Get(k K) (V, bool) {
	val, exists := m.input[k]
	return val, exists
}

func (m *MapOperation[K, V]) GetOrDefault(k K, defaultVal V) V {
	if val, exists := m.input[k]; exists {
		return val
	} else {
		return defaultVal
	}

}

func (m *MapOperation[K, V]) Put(k K, v V) (V, bool) {
	if val, exists := m.input[k]; exists {
		m.input[k] = v
		return val, true
	} else {
		m.input[k] = v
		return val, false
	}
}

func (m *MapOperation[K, V]) PutIfAbsent(k K, v V) (V, bool) {
	if val, exists := m.input[k]; exists {
		return val, true
	} else {
		m.input[k] = v
		return val, false
	}
}

func (m *MapOperation[K, V]) PutOrDefault(k K, defaultVal V) V {
	if val, exists := m.input[k]; exists {
		return val
	} else {
		m.input[k] = defaultVal
		return defaultVal
	}
}

func (m *MapOperation[K, V]) Replace(k K, v V) (V, bool) {
	if val, exists := m.input[k]; exists {
		m.input[k] = v
		return val, true
	} else {
		return val, false
	}
}

func (m *MapOperation[K, V]) ReplaceWith(k K, v1 V, v2 V) bool {
	if val, exists := m.input[k]; exists {
		if any(val) == any(v1) {
			m.input[k] = v2
			return true
		}
		// 值不相等 也是false
		return false
	} else {
		return false
	}
}

func (m *MapOperation[K, V]) PutAll(m2 map[K]V) {
	for k, v := range m2 {
		m.input[k] = v
	}
}

func (m *MapOperation[K, V]) Compute(k K, f func(K, V) V) (V, bool) {
	if val, exists := m.input[k]; exists {
		newValue := f(k, val)
		m.input[k] = newValue
		return newValue, true
	} else {
		newValue := f(k, val)
		m.input[k] = newValue
		return newValue, false
	}
}

func (m *MapOperation[K, V]) ComputeIfAbsent(k K, f func(K) V) (V, bool) {
	if val, exists := m.input[k]; !exists {
		newValue := f(k)
		m.input[k] = newValue
		return newValue, false
	} else {
		return val, true
	}
}

func (m *MapOperation[K, V]) ComputeIfPresent(k K, f func(K, V) V) (V, bool) {
	if val, exists := m.input[k]; !exists {
		return val, false
	} else {
		newValue := f(k, val)
		m.input[k] = newValue
		return newValue, true
	}
}

func (m *MapOperation[K, V]) Merge(k K, val V, f func(oldVal V, val V) V) (V, bool) {
	if oldVal, exists := m.input[k]; !exists {
		m.input[k] = val
		return val, false
	} else {
		newValue := f(oldVal, val)
		m.input[k] = newValue
		return newValue, true
	}
}

func (m *MapOperation[K, V]) Remove(k K) V {
	val := m.input[k]
	delete(m.input, k)
	return val
}

func (m *MapOperation[K, V]) RemoveValue(k K, v V) bool {
	if val, exists := m.input[k]; exists {
		if any(val) == any(v) {
			delete(m.input, k)
			return true
		}
		return false
	} else {
		return false
	}
}

func (m *MapOperation[K, V]) Keys() *[]K {
	i := 0
	r := make([]K, len(m.input))
	for k := range m.input {
		r[i] = k
		i++
	}
	return &r
}

func (m *MapOperation[K, V]) Values() *[]V {
	i := 0
	r := make([]V, len(m.input))
	for _, v := range m.input {
		r[i] = v
		i++
	}
	return &r
}

func (m *MapOperation[K, V]) ContainsKey(k K) bool {
	_, exists := m.input[k]
	return exists
}

func (m *MapOperation[K, V]) ContainsValue(v1 V) bool {
	for _, v := range m.input {
		if any(v) == any(v1) {
			return true
		}
	}
	return false
}

func (m *MapOperation[K, V]) String() string {
	return fmt.Sprintf("%v", m.input)
}

func (m *MapOperation[K, V]) ToJsonString() (string, error) {
	marshal, err := json.Marshal(m.input)
	if err != nil {
		return "", err
	}
	return convert_helper.BytesToStr(marshal), nil
}

func (m *MapOperation[K, V]) Size() int {
	return len(m.input)
}

func (m *MapOperation[K, V]) Copy() map[K]V {
	r := make(map[K]V, len(m.input))
	for k, v := range m.input {
		r[k] = v
	}
	return r
}

func (m *MapOperation[K, V]) Clear() {
	for k := range m.input {
		delete(m.input, k)
	}
}

func (m *MapOperation[K, V]) Update(m2 map[K]V) {
	for k, v := range m2 {
		m.input[k] = v
	}
}

func (m *MapOperation[K, V]) KeyMapSet() mapset.Set[K] {
	sets := mapset.NewSet[K]()
	for k := range m.input {
		sets.Add(k)
	}
	return sets
}

func (m *MapOperation[K, V]) KeyThreadSafeMapSet() mapset.Set[K] {
	sets := mapset.NewThreadUnsafeSet[K]()
	for k := range m.input {
		sets.Add(k)
	}
	return sets
}

func (m *MapOperation[K, V]) KeyHashSet() *hashset.Set[K] {
	sets := hashset.New[K]()
	for k := range m.input {
		sets.Add(k)
	}
	return sets
}

func (m *MapOperation[K, V]) KeyLinkedHashSet() *linkedhashset.Set[K] {
	sets := linkedhashset.New[K]()
	for k := range m.input {
		sets.Add(k)
	}
	return sets
}

func (m *MapOperation[K, V]) KeyTreeSet(comparator utils.Comparator[K]) *treeset.Set[K] {
	sets := treeset.NewWith[K](comparator)
	for k := range m.input {
		sets.Add(k)
	}
	return sets
}

func (m *MapOperation[K, V]) Each(f func(K, V) bool) {
	for k, v := range m.input {
		if f(k, v) {
			break
		}
	}
}

func (m *MapOperation[K, V]) Iter() <-chan *map_helper.Pair[K, V] {
	return map_helper.Iter[K, V](m.input)
}

func (m *MapOperation[K, V]) GetIterator() *map_helper.Iterator[*map_helper.Pair[K, V]] {
	return map_helper.GetIterator[K, V](m.input)
}

func NewEmptyMapOperation[K comparable, V any]() *MapOperation[K, V] {
	return &MapOperation[K, V]{
		input: map[K]V{},
	}
}

func NewMapOperation[K comparable, V any](input map[K]V) *MapOperation[K, V] {
	return &MapOperation[K, V]{
		input: input,
	}
}

func NewMapFromKeys[K comparable, V any](keys *[]K) *MapOperation[K, V] {
	return &MapOperation[K, V]{
		input: map_helper.FromKeys[K, V](keys),
	}
}

func NewMapFromKVS[K comparable, V any](keys *[]K, values *[]V) *MapOperation[K, V] {
	return &MapOperation[K, V]{
		input: map_helper.FromKVs[K, V](keys, values),
	}
}

func NewMapFromPairFunc[K comparable, S any, V any](source []S, fn func(S) (K, V)) *MapOperation[K, V] {
	return &MapOperation[K, V]{
		input: map_helper.FromPairFunc[K, S, V](source, fn),
	}
}
