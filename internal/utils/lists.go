package utils

import "math/rand"

func ListRandomElement[T any](list []T) T {
	return list[rand.Intn(len(list))]
}
