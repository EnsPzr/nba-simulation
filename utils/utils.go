package utils

type number interface {
	int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64 |
		int | uint | float32 | float64
}

func Remove[X number | bool | string](arr []X, deleted X) []X {
	index := FindIndex(arr, deleted)
	if index == -1 {
		return arr
	}
	return append(arr[:index], arr[index+1:]...)
}

func FindIndex[X number | bool | string](arr []X, search X) int {
	for index, val := range arr {
		if val == search {
			return index
		}
	}
	return -1
}
