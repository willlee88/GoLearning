package main

import "fmt"

func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func Filter[T any](s []T, keep func(T) bool) []T {
	out := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func main() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println(Keys(m))
	fmt.Println(Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }))
}
