package arraysx

// RemoveByIndex 移除列表索引处的元素
func removeByIndex[T any](arr []T, i int) (res []T) {
	switch i {
	case 0:
		res = arr[1:]
	case len(arr) - 1:
		res = arr[:i]
	default:
		// 注: `append(arr[:i], arr[i+1:]...)`是错误用法，将导致原数组从`i索引起的元素`被`arr[i+1:]`覆盖
		res = append(res, arr[:i]...)
		res = append(res, arr[i+1:]...)
	}
	return
}

// RemoveByIndex 移除列表索引处的元素，支持负索引，如: -1
func RemoveByIndex[T any](arr []T, i int) (res []T) {
	if i < 0 {
		return removeByIndex(arr, len(arr)+i)
	}
	return removeByIndex(arr, i)
}

// RemoveByCondition 移除列表中与条件相匹配的所有元素
func RemoveByCondition[T any](arr []T, condition func(T) bool) (res []T) {
	for i := range arr {
		if !condition(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return
}

// RemoveRedundant 移除重复的多余的元素(重复元素仅保留最后一个)
func RemoveRedundant[T any](arr []T, isEqual func(T, T) bool) (res []T) {
	for i, l, r := 0, len(arr), false; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if isEqual(arr[i], arr[j]) {
				r = true
				break
			}
		}
		if !r {
			res = append(res, arr[i])
		}
		r = false
	}
	return
}

// Index 检索列表中与条件相匹配的元素位置
func Index[T any](arr []T, condition func(T) bool) int {
	for i := 0; i < len(arr); i++ {
		if condition(arr[i]) {
			return i
		}
	}
	return -1
}

func Filter[T any](arr []T, cp func(T) bool) (t []T) {
	for i := 0; i < len(arr); i++ {
		if cp(arr[i]) {
			t = append(t, arr[i])
		}
	}
	return
}

func Split[T any](arr []T, cp func(T) bool) (t, f []T) {
	for i := 0; i < len(arr); i++ {
		if cp(arr[i]) {
			t = append(t, arr[i])
		} else {
			f = append(f, arr[i])
		}
	}
	return
}
