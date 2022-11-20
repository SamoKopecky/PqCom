package dilithium

import "fmt"

func KeyGen() {
	zeta := genRand(256)
	hash := shake256(zeta, 128)
	ro := hash[:32]
	ro_dash := hash[32:96]
	// K := hash[96:]

	A_hat := expandA(ro)

	s := expandS(ro_dash)
	s1 := s[:L]
	s2 := s[L:]

	fmt.Printf("Before %d \n\n", A_hat[0][0])
	ntt(A_hat[0][0])
	invNtt(A_hat[0][0])
	fmt.Printf("After %d", A_hat[0][0])
	fmt.Printf("\n\nA: %d, s1: %d, s2:%d", A_hat[0][0][0], s1[0][0], s2[0][0])

}
