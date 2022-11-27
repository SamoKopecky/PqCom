package dilithium

func KeyGen() (pk, sk []byte) {
	zeta := shake256(genRand(256), 128)
	rho := zeta[:32]
	rho_dash := zeta[32:96]
	K := zeta[96:]

	A_hat := expandA(rho)

	s := expandS(rho_dash)
	s_1 := s[:L]
	s_2 := s[L:]

	s_1_hat := nttPolyVec(s_1)
	t := make([][]int, L)
	for i := 0; i < L; i++ {
		t[i] = mulPolyVec(A_hat[i], s_1_hat)
	}
	t = addPolyVec(invNttPolyVec(t), s_2)

	t_1, t_0 := powerToModPolyVec(t, D)

	t_1_packed := bitPackPolyVec(t_1, 10)
	shake := append(rho, t_1_packed...)
	tr := shake256(shake, 32)

	pk = append(rho, t_1_packed...)

	sk = append(rho, K...)
	sk = append(sk, tr...)
	sk = append(sk, bitPackAlteredPolyVec(s_1, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(s_2, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(t_0, 1<<13, 13)...)
	return
}

func Sign(sk []byte, message []byte) (sigma []byte) {
	rho := sk[:32]
	K := sk[32:64]
	tr := sk[64:96]
	s_1 := modPMPolyVec(bitUnpackAlteredPolyVec(sk[96:SBytes+96], Eta, 3), Q)
	s_2 := modPMPolyVec(bitUnpackAlteredPolyVec(sk[96+SBytes:SBytes*2+96], Eta, 3), Q)
	t_0 := modPMPolyVec(bitUnpackAlteredPolyVec(sk[96+SBytes*2:], 1<<13, 13), 1<<D)
	
	A_hat := expandA(rho)

	s_1_hat := nttPolyVec(s_1)
	s_2_hat := nttPolyVec(s_2)
	t_0_hat := nttPolyVec(t_0)

	shake := append(tr, message...)
	mi := shake256(shake, 64)

	shake = append(K, mi...)
	rho_dash := shake256(shake, 64)
	kappa := -L

	for true {
		kappa += L

		y := expandMask(rho_dash, kappa)
		y_hat := nttPolyVec(y)

		w := make([][]int, L)
		for i := 0; i < L; i++ {
			w[i] = invNtt(mulPolyVec(A_hat[i], y_hat))
		}
		w_1 := highBitsPolyVec(w, 2*GammaTwo)

		w_1_packed := bitPackPolyVec(w_1, 6)

		shake = append(mi, w_1_packed...)
		c_wave := shake256(shake, 32)

		c := sampleInBall(c_wave)
		c_hat := ntt(c)

		cs_1 := invNttPolyVec(scalePolyVecByPoly(s_1_hat, c_hat))
		cs_2 := invNttPolyVec(scalePolyVecByPoly(s_2_hat, c_hat))

		z := modPMPolyVec(addPolyVec(y, cs_1), Q)

		w_sub_cs_2 := modPMPolyVec(subPolyVec(w, cs_2), Q)

		r_0 := lowBitsPolyVec(modPMPolyVec(w_sub_cs_2, Q), 2*GammaTwo)

		if checkNormPolyVec(z, GammaOne-Beta) || checkNormPolyVec(r_0, GammaTwo-Beta) {
			continue
		}

		ct_0 := modPMPolyVec(invNttPolyVec(scalePolyVecByPoly(t_0_hat, c_hat)), Q)
		ct_0_inv := modPMPolyVec(inversePolyVec(ct_0), Q)
		w_sub_cs_2_add_ct_0 := modPMPolyVec(addPolyVec(w_sub_cs_2, ct_0), Q)

		h := makeHintPolyVec(ct_0_inv, w_sub_cs_2_add_ct_0, 2*GammaTwo)

		ones := 0
		for i := 0; i < len(h); i++ {
			for j := 0; j < N; j++ {
				if h[i][j] == 1 {
					ones++
				}
			}
		}
		if checkNormPolyVec(modPMPolyVec(ct_0, Q), GammaTwo) {
			continue
		}
		if ones > Omega {
			continue
		}

		z_packed := bitPackAlteredPolyVec(z, GammaOne, 18)
		h_packed := bitPackHint(h)
		sigma = append(c_wave, z_packed...)
		sigma = append(sigma, h_packed...)
		break
	}
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

	shake := append(rho, t_1_bytes...)
	shake = append(shake256(shake, 32), message...)
	mi := shake256(shake, 64)

	c := sampleInBall(c_wave)
	c_hat := ntt(c)

	z_hat := nttPolyVec(z)

	Az := make([][]int, L)
	for i := 0; i < L; i++ {
		Az[i] = mulPolyVec(A_hat[i], z_hat)
	}

	t_1_2_to_d := modPMPolyVec(scalePolyVecByInt(t_1, 1<<D), Q)
	t_1_2_to_d_hat := nttPolyVec(t_1_2_to_d)
	ct_1 := scalePolyVecByPoly(t_1_2_to_d_hat, c_hat)

	r := invNttPolyVec(subPolyVec(Az, ct_1))

	w_1 := useHintPolyVec(h, modPMPolyVec(r, Q), 2*GammaTwo)
	shake = append(mi, bitPackPolyVec(w_1, 6)...)
	verified = BytesEqual(c_wave, shake256(shake, 32))
	return
}
