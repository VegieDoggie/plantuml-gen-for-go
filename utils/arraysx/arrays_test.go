package arraysx

import (
	"fmt"
	"testing"
)

func TestRemoveByCondition(t *testing.T) {
	x := []int{1, 2, 3, 4, 5}
	fmt.Println(x)

	fmt.Println(RemoveByCondition(x, func(a int) bool {
		return a == 2
	}))
	fmt.Println(RemoveByCondition(x, func(a int) bool {
		return a > 3
	}))
}

func TestRemoveByIndex(t *testing.T) {
	x := []int{1, 2, 3, 4, 5}
	fmt.Println(x)

	fmt.Println(RemoveByIndex(x, 1))
	fmt.Println(RemoveByIndex(x, 0))
}

func TestRemoveAt(t *testing.T) {
	x := []int{1, 2, 3, 4, 5}
	fmt.Println(x)

	fmt.Println(RemoveByIndex(x, -2))
	fmt.Println(RemoveByIndex(x, -1))
}

func TestUnrepeatable(t *testing.T) {
	x := []int{1, 2, 2, 3, 3, 4, 5, 1}
	fmt.Println(x)
	fmt.Println(RemoveRedundant(x, func(a, b int) bool { return a == b }))
}
