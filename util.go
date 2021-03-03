package devmapper

func roundUp(n int, divider int) int {
	return (n + divider - 1) / divider * divider
}
