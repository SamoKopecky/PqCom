package kyber

func (kyb *kyber) genRow(rho []byte, a, b byte) []int {
	return kyb.parse(xof(rho, a, b, n*3))
}

func (kyb *kyber) genTRow(rho []byte, a, b byte) []int {
	return kyb.parse(xof(rho, b, a, n*3))
}

func (kyb *kyber) genPolyMat(rho []byte, transpose bool) (mat [][][]int) {
	generate := kyb.genRow
	if transpose {
		generate = kyb.genTRow
	}
	for i := byte(0); i < byte(kyb.k); i++ {
		row := [][]int{}
		for j := byte(0); j < byte(kyb.k); j++ {
			row = append(row, generate(rho, j, i))
		}
		mat = append(mat, row)
	}
	return
}
