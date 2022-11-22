package dilithium

import "fmt"

func KeyGen() (pk, sk []byte) {
	zeta := genRand(256)
	hash := shake256(zeta, 128)
	rho := hash[:32]
	rho_dash := hash[32:96]
	K := hash[96:]

	A_hat := expandA(rho)

	s := expandS(rho_dash)
	s_1 := s[:L]
	s_2 := s[L:]
	s_1_hat := nttPolyVec(s_1)

	t := make([][]int, L)
	for i := 0; i < L; i++ {
		t[i] = mulPolyVec(A_hat[i], s_1_hat)
	}
	t = reducePolyVec(addPolyVec(invNttPolyVec(t), s_2))

	t_1, t_0 := powerToModPolyVec(t, D)

	t_1_packed := bitPackPolyVec(t_1, 10)
	shake_input := rho
	shake_input = append(shake_input, t_1_packed...)
	tr := shake256(shake_input, 32)
	pk = rho
	pk = append(pk, t_1_packed...)
	sk = rho
	sk = append(sk, K...)
	sk = append(sk, tr...)
	sk = append(sk, bitPackAlteredPolyVec(s_1, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(s_2, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(t_0, 1<<13, 13)...)
	fmt.Printf("%d\n", t_0[0][:20])
	return
}

func Sign(sk []byte, message []byte) (sigma []byte) {
	rho := sk[:32]
	K := sk[32:64]
	tr := sk[64:96]
	s_1 := reducePolyVec(bitUnpackAlteredPolyVec(sk[96:SBytes+96], Eta, 3)) 
	s_2 := reducePolyVec(bitUnpackAlteredPolyVec(sk[96+SBytes:SBytes*2+96], Eta, 3)) 
	t_0 := reducePolyVec(bitUnpackAlteredPolyVec(sk[96+SBytes*2:], 1<<13, 13)) 
	fmt.Printf("%d\n", reducePolyVec(t_0)[0][:20])
	A_hat := expandA(rho)

	s_1_hat := nttPolyVec(s_1)
	s_2_hat := nttPolyVec(s_2)
	t_0_hat := nttPolyVec(t_0)

	shake_input := tr
	shake_input = append(shake_input, message...)
	mi := shake256(shake_input, 64)

	shake_input = K
	shake_input = append(shake_input, mi...)
	rho_dash := shake256(shake_input, 64)
	kappa := -L

	omegas := 0
	norms := 0
	for true {
		kappa += L

		y := expandMask(rho_dash, kappa)
		y_hat := nttPolyVec(y)
		w := make([][]int, L)
		for i := 0; i < L; i++ {
			w[i] = invNtt(mulPolyVec(A_hat[i], y_hat))
		}
		w = reducePolyVec(w)
		w_1 := highBitsPolyVec(w, 2*GammaTwo)

		w_1_packed := bitPackPolyVec(w_1, 6)

		shake_input = mi
		shake_input = append(shake_input, w_1_packed...)
		c_wave := shake256(shake_input, 32)

		c := sampleInBall(c_wave)
		c_hat := ntt(c)

		c_times_s_1 := invNttPolyVec(scalePolyVecByPoly(s_1_hat, c_hat))
		c_times_s_2 := reducePolyVec(invNttPolyVec(scalePolyVecByPoly(s_2_hat, c_hat)))
		c_times_t_0 := reducePolyVec(invNttPolyVec(scalePolyVecByPoly(t_0_hat, c_hat)))

		z := reducePolyVec(addPolyVec(y, reducePolyVec(c_times_s_1)))

		w_minus_c_times_s_2 := reducePolyVec(subPolyVec(w, c_times_s_2))

		r_0 := lowBitsPolyVec(w_minus_c_times_s_2, 2*GammaTwo)

		if checkNormPolyVec(z, GammaOne-Beta) || checkNormPolyVec(reducePolyVec(r_0), GammaTwo-Beta) {
			norms++
			continue
		}
		c_times_t_0_reduced := reducePolyVec(inversePolyVec(c_times_t_0))
		second := reducePolyVec(addPolyVec(w_minus_c_times_s_2, c_times_t_0))
		h := makeHintPolyVec(c_times_t_0_reduced, second, 2*GammaTwo)
		// This works
		test := useHintPolyVec(h, second, 2*GammaTwo)
		// test2 := highBitsPolyVec(addPolyVec(second, c_times_t_0_reduced), 2*GammaTwo)
		// fmt.Printf("%d\n", test)
		// fmt.Printf("%d\n", test2)

		ones := 0
		for i := 0; i < len(h); i++ {
			for j := 0; j < N; j++ {
				if h[i][j] == 1 {
					ones++
				}
			}
		}
		if checkNormPolyVec(reducePolyVec(c_times_t_0), GammaTwo) {
			norms++
			continue
		}
		if ones > Omega {
			omegas++
			continue
		}
		z_packed := bitPackAlteredPolyVec(z, GammaOne, 18)
		h_packed := bitPackHint(h)
		sigma = append(sigma, c_wave...)
		sigma = append(sigma, z_packed...)
		sigma = append(sigma, h_packed...)
		fmt.Printf("w_1: %d\n", w_1[0][:50])
		fmt.Printf("used %d\n", test[0][:50])
		break
	}
	// fmt.Printf("\nKappa: %d", kappa/L)
	// fmt.Printf("\nOmegas: %d", omegas)
	// fmt.Printf("\nNorms: %d", norms)
	return
}

func Verify(pk, message, sigma []byte) (verified bool) {
	rho := pk[:32]
	t_1_bytes := pk[32:]
	c_wave := sigma[:32]
	z := bitUnpackAlteredPolyVec(sigma[32:ZBytes+32], GammaOne, 18)
	h := bitUnpackHint(sigma[ZBytes+32:])
	t_1 := bitUnpackPolyVec(t_1_bytes, 10)

	A_hat := expandA(rho)

	shake_b := rho
	shake_b = append(shake_b, t_1_bytes...)
	mi := shake256(shake_b, 64)
	shake_b = mi
	shake_b = append(shake_b, message...)
	mi = shake256(shake_b, 64)

	c := sampleInBall(c_wave)

	c_hat := ntt(c)
	z_hat := nttPolyVec(z)

	A_times_z := make([][]int, L)
	for i := 0; i < L; i++ {
		A_times_z[i] = mulPolyVec(A_hat[i], z_hat)
	}

	t_1_times_2_d := scalePolyVecByInt(t_1, 1<<D)
	t_1_times_2_d_hat := nttPolyVec(t_1_times_2_d)
	c_times_t_1 := scalePolyVecByPoly(t_1_times_2_d_hat, c_hat)

	r := invNttPolyVec(subPolyVec(A_times_z, c_times_t_1))

	fmt.Printf("w_1: %d\n", useHintPolyVec(h, r, 2*GammaTwo)[0][:50])
	w_1_dash := useHintPolyVec(h, r, 2*GammaTwo)
	fmt.Printf("w_1: %d\n", w_1_dash[0][:50])
	// fmt.Printf("hhh: %02d\n", h[0][:50])
	// fmt.Printf("rrr: %d\n", highBitsPolyVec(r, GammaTwo*2)[0][:50])

	shake_b = mi
	shake_b = append(shake_b, bitPackPolyVec(w_1_dash, 6)...)
	other_c_wave := shake256(shake_b, 32)
	fmt.Printf("%d\n", c_wave[:10])
	fmt.Printf("%d\n", other_c_wave[:10])
	verified = BytesEqual(c_wave, other_c_wave)
	return
}
