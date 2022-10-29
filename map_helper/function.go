package map_helper

// FromKeys 用 keys 来造map
func FromKeys[K comparable, V any](keys *[]K) map[K]V {
	raw := *keys
	if len(raw) <= 0 {
		return map[K]V{}
	}
	res := make(map[K]V, len(raw))

	for _, k := range raw {
		var v V
		res[k] = v
	}
	return res
}

// FromKVs 用  keys 和 vals 来造map len(keys) == len(vals)
func FromKVs[K comparable, V any](keys *[]K, vals *[]V) map[K]V {
	raw := *keys
	rawLength := len(raw)
	if rawLength <= 0 {
		return map[K]V{}
	}
	raw2 := *vals
	raw2Length := len(raw2)
	if raw2Length <= 0 {
		return map[K]V{}
	}
	if rawLength != raw2Length {
		panic("keys_length is not equal to vals_length")
	}
	res := make(map[K]V, rawLength)
	for i := 0; i < rawLength; i++ {
		res[raw[i]] = raw2[i]
	}
	return res
}

// FromPairArr 用二维数组来造 kv
//func FromPairArr[K comparable, V any](kv *[][]any) map[K]V {
//	raw := *kv
//	rawLength := len(raw)
//	if rawLength <= 0 {
//		return map[K]V{}
//	}
//	res := make(map[K]V, rawLength)
//	for _, item := range raw {
//		if len(item) < 2 {
//			panic("kv length must >=2 ")
//		}
//		res[item[0]] = item[1]
//	}
//	return res
//}

// FromPairFunc 用返回两个值的函数来造kv
// K 是 key S 是 提供值的容器 V 是 值
func FromPairFunc[K comparable, S any, V any](source []S, fn func(S) (K, V)) map[K]V {
	rawLength := len(source)
	if rawLength <= 0 {
		return map[K]V{}
	}
	res := make(map[K]V, rawLength)
	for _, item := range source {
		k, v := fn(item)
		res[k] = v
	}
	return res
}
