package sliceutils

func Sum(slice []int) int {
	sum := 0
	for i := range slice {
		sum += slice[i]
	}
	return sum
}

func Min(slice []int) int {
	min := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
		}
	}
	return min
}

func Max(slice []int) int {
	max := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}

func ContainsString(slice []string, val string) bool {
	for i := range slice {
		if slice[i] == val {
			return true
		}
	}
	return false
}

func IntersectString(slice1, slice2 []string) []string {
	intersection := make([]string, 0)
	for i := range slice1 {
		if ContainsString(slice2, slice1[i]) {
			intersection = append(intersection, slice1[i])
		}
	}
	return intersection
}

func DistinctString(slice []string) []string {
	newSlice := make([]string, 0, len(slice))
	for i := range slice {
		if ContainsString(newSlice, slice[i]) {
			continue
		}
		newSlice = append(newSlice, slice[i])
	}
	return newSlice
}

func Reverse(slice []byte) {
	for j, k := 0, len(slice)-1; j < len(slice)/2; j, k = j+1, k-1 {
		slice[j], slice[k] = slice[k], slice[j]
	}
}

func GridString(grid [][]byte) string {
	s := ""
	for row := range grid {
		s += "\n" + string(grid[row])
	}
	return s[1:]
}
