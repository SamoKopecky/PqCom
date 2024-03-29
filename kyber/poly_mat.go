package kyber

func (kyb *Kyber) genRow(rho []byte, a, b byte) []int {
	return kyb.parse(kyb.xof(rho, a, b, n*3))
}

func (kyb *Kyber) genTRow(rho []byte, a, b byte) []int {
	return kyb.parse(kyb.xof(rho, b, a, n*3))
}

func (kyb *Kyber) genPolyMat(rho []byte, transpose bool) (mat [][][]int) {
	generate := kyb.genRow
	if transpose {
		generate = kyb.genTRow
	}
	mat = make([][][]int, kyb.k)
	for i := byte(0); i < byte(kyb.k); i++ {
		mat[i] = make([][]int, kyb.k)
		for j := byte(0); j < byte(kyb.k); j++ {
			mat[i][j] = generate(rho, j, i)
		}
	}
	return
}
