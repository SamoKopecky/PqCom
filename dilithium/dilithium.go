package dilithium

import "fmt"

func KeyGen() (pk, sk []byte) {
	zeta := genRand(256)
	hash := shake256(zeta, 128)
	ro := hash[:32]
	ro_dash := hash[32:96]
	K := hash[96:]

	A_hat := expandA(ro)

	s := expandS(ro_dash)
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
	shake_input := ro
	shake_input = append(shake_input, t_1_packed...)
	tr := shake256(shake_input, 32)

	pk = ro
	pk = append(pk, t_1_packed...)
	sk = ro
	sk = append(sk, K...)
	sk = append(sk, tr...)
	sk = append(sk, bitPackAlteredPolyVec(s_1, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(s_2, Eta, 3)...)
	sk = append(sk, bitPackAlteredPolyVec(t_0, 1<<13, 13)...)
	return
}

func Sign(sk []byte, message []byte) (sigma []byte) {
	ro := sk[:32]
	K := sk[32:64]
	tr := sk[64:96]
	s_1 := bitUnpackAlteredPolyVec(sk[96:S_BYTES+96], Eta, 3)
	s_2 := bitUnpackAlteredPolyVec(sk[96+S_BYTES:S_BYTES*2+96], Eta, 3)
	t_0 := bitUnpackAlteredPolyVec(sk[96+S_BYTES*2:], 1<<13, 13)
	A_hat := expandA(ro)

	s_1_hat := nttPolyVec(s_1)
	s_2_hat := nttPolyVec(s_2)
	t_0_hat := nttPolyVec(t_0)

	shake_input := tr
	shake_input = append(shake_input, message...)
	mi := shake256(shake_input, 64)
	shake_input = K
	shake_input = append(shake_input, mi...)
	ro_dash := shake256(shake_input, 64)
	kappa := -L

	omegas := 0
	norms := 0
	for true {
		kappa += L

		y := expandMask(ro_dash, kappa)

		w := make([][]int, L)
		for i := 0; i < L; i++ {
			w[i] = invNtt(mulPolyVec(A_hat[i], nttPolyVec(y)))
		}
		w = reducePolyVec(w)
		w_1 := highBitsPolyVec(w, 2*GammaTwo)
		w_1_packed := bitPackPolyVec(w_1, 6)

		shake_input = mi
		shake_input = append(shake_input, w_1_packed...)
		c_wave := shake256(shake_input, 32)
		c := sampleInBall(c_wave)
		c_hat := ntt(c)

		c_times_s_1 := invNttPolyVec(scalePolyVec(s_1_hat, c_hat))
		c_times_s_2 := invNttPolyVec(scalePolyVec(s_2_hat, c_hat))
		c_times_t_0 := invNttPolyVec(scalePolyVec(t_0_hat, c_hat))

		z := addPolyVec(y, c_times_s_1)

		w_minus_c_times_s_2 := subPolyVec(w, reducePolyVec(c_times_s_2))

		r_0 := lowBitsPolyVec(w_minus_c_times_s_2, 2*GammaTwo)

		if checkNormPolyVec(reducePolyVec(z), GammaOne-Beta) || checkNormPolyVec(reducePolyVec(r_0), GammaTwo-Beta) {
			norms++
			continue
		}

		second := subPolyVec(w, c_times_s_2)
		second = addPolyVec(second, c_times_t_0)
		h := makeHintPolyVec(inversePolyVec(c_times_t_0), second, 2*GammaTwo)
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
		break
	}
	fmt.Printf("\nKappa: %d", kappa/L)
	fmt.Printf("\nOmegas: %d", omegas)
	fmt.Printf("\nNorms: %d", norms)
	return
}
