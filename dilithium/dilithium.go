package dilithium

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
	t = addPolyVec(invNttPolyVec(t), s_2)

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

func Sign(pk, sk []byte, message []byte) (c_wave, z, h []byte) {
	ro := pk[:32]
	t_1 := bitUnpackPolyVec(pk[32:], 10)
	A_hat := expandA(ro)

	K := sk[32:64]
	tr := sk[64:96]
	s_1 := bitUnpackAlteredPolyVec(sk[96:S_BYTES+96], Eta, 3)
	s_2 := bitUnpackAlteredPolyVec(sk[96+S_BYTES:S_BYTES*2+96], Eta, 3)
	t_0 := bitUnpackAlteredPolyVec(sk[96+S_BYTES*2:], 1<<13, 13)

	shake_input := tr
	shake_input = append(shake_input, message...)
	mi := shake256(shake_input, 64)
	kappa := 0

	shake_input = K
	shake_input = append(shake_input, mi...)
	ro_dash := shake256(shake_input, 64)

	return
}
