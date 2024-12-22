package graph

func reversePermutation(p []uint) []uint {
	n := len(p)
	q := make([]uint, n)
	for i, v := range p {
		q[v] = uint(i)
	}
	return q
}
