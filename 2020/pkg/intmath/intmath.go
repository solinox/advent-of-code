package intmath

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(ints ...int) int {
	if len(ints) < 2 {
		ints = append([]int{1, 1}, ints...)
		return LCM(ints...)
	}

	a, b := ints[0], ints[1]
	ints = ints[2:]

	result := a * b / GCD(a, b)

	for i := 0; i < len(ints); i++ {
		result = LCM(result, ints[i])
	}

	return result
}
