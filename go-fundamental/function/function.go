package function

func Add(a int, b int) int {
	return a + b
}

func Multiply(arr ...int) int {
	ans := 1
	for _, val := range arr {
		ans *= val
	}
	return ans
}
