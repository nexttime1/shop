package main

func Add[T int32 | int64 | uint | float32 | float64](x, y T) T {
	return x + y
}

type MyMap[key int32 | float32 | float64, value string | float64] map[key]value

type Man struct {
}
type Woman struct {
}

type Company[T Man | Woman] struct {
	Name   string
	Person T
}

type MyChan[T int | string] chan T
