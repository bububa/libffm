package utils

import (
	"math"
	"math/rand"

	"github.com/bububa/libffm"
)

func Randomization(l int, randFlag bool) []int {
	order := make([]int, l)
	for i := 0; i < len(order); i++ {
		order[i] = i
	}
	if randFlag {
		for i := len(order); i > 1; i-- {
			tmp := order[i-1]
			index := rand.Intn(i)
			order[i-1] = order[index]
			order[index] = tmp
		}
	}
	return order
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func GetLatentFactorsNumberAligned(latentFactorsNumber int) int {
	return int(math.Ceil(float64(latentFactorsNumber)/libffm.ALIGN) * libffm.ALIGN)
}
