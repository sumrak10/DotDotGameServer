package utils

func MapToIndexedKeysMap[K comparable, V any](originalMap map[K]V) map[uint]K {
	result := make(map[uint]K)
	var i uint
	for key := range originalMap {
		result[i] = key
		i++
	}
	return result
}

func MapToIndexedValuesMap[K comparable, V any](originalMap map[K]V) map[uint]V {
	result := make(map[uint]V)
	var i uint
	for _, val := range originalMap {
		result[i] = val
		i++
	}
	return result
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func MapValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
