package kyber

func genRow(rho []byte, a, b byte) []int {
	return parse(xof(rho, a, b, N*3))
}

func genTRow(rho []byte, a, b byte) []int {
	return parse(xof(rho, b, a, N*3))
}

func genPolyMat(rho []byte, transpose bool) (mat [][][]int) {
	generate := genRow
	if transpose {
		generate = genTRow
	}
	for i := byte(0); i < byte(K); i++ {
		row := [][]int{}
		for j := byte(0); j < byte(K); j++ {
			row = append(row, generate(rho, j, i))
		}
		mat = append(mat, row)
	}
	return
}
