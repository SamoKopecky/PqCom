package dilithium

import (
	"github.com/SamoKopecky/pqcom/main/myio"
)

const (
	q = 8380417
	d = 13
	n = 256
)

type Dilithium struct {
	tau      int
	k        int
	l        int
	eta      int
	s1Bytes  int
	s2Bytes  int
	zBytes   int
	gammaOne int
	gammaTwo int
	omega    int
	beta     int
	PkSize   int
	SkSize   int
	SigSize  int
	zBits    int
	w_1Bits  int
	sBits    int
}

func Dilithium2() Dilithium {
	dil := Dilithium{
		tau:      39,
		k:        4,
		l:        4,
		eta:      2,
		gammaOne: 1 << 17,
		omega:    80,
		PkSize:   1312,
		SkSize:   2528,
		SigSize:  2420,
		zBits:    18,
		w_1Bits:  6,
		sBits:    3,
	}
	dil.s1Bytes = (n * dil.sBits / 8) * dil.l
	dil.s2Bytes = (n * dil.sBits / 8) * dil.k
	dil.zBytes = (n * dil.zBits / 8) * dil.l
	dil.gammaTwo = (q - 1) / 88
	dil.beta = dil.tau * dil.eta
	return dil
}

func Dilithium3() Dilithium {
	dil := Dilithium{
		tau:      49,
		k:        6,
		l:        5,
		eta:      4,
		gammaOne: 1 << 19,
		omega:    55,
		PkSize:   1952,
		SkSize:   4000,
		SigSize:  3293,
		zBits:    20,
		w_1Bits:  4,
		sBits:    4,
	}
	dil.s1Bytes = (n * dil.sBits / 8) * dil.l
	dil.s2Bytes = (n * dil.sBits / 8) * dil.k
	dil.zBytes = (n * dil.zBits / 8) * dil.l
	dil.gammaTwo = (q - 1) / 32
	dil.beta = dil.tau * dil.eta
	return dil
}

func Dilithium5() Dilithium {
	dil := Dilithium{
		tau:      60,
		k:        8,
		l:        7,
		eta:      2,
		gammaOne: 1 << 19,
		omega:    75,
		PkSize:   2592,
		SkSize:   4864,
		SigSize:  4595,
		zBits:    20,
		w_1Bits:  4,
		sBits:    3,
	}
	dil.s1Bytes = (n * dil.sBits / 8) * dil.l
	dil.s2Bytes = (n * dil.sBits / 8) * dil.k
	dil.zBytes = (n * dil.zBits / 8) * dil.l
	dil.gammaTwo = (q - 1) / 32
	dil.beta = dil.tau * dil.eta
	return dil
}

func (dil *Dilithium) KeyGen() (pk, sk []byte) {
	zeta := dil.shake256(dil.genRand(n), 128)
	rho := zeta[:32]
	rho_dash := zeta[32:96]
	K := zeta[96:]

	A_hat := dil.expandA(rho)

	s := dil.expandS(rho_dash)
	s_1 := s[:dil.l]
	s_2 := s[dil.l:]

	s_1_hat := dil.nttPolyVec(s_1)
	t := make([][]int, dil.k)
	for i := 0; i < dil.k; i++ {
		t[i] = dil.mulPolyVec(A_hat[i], s_1_hat)
	}
	t = dil.addPolyVec(dil.invNttPolyVec(t), s_2)

	t_1, t_0 := dil.powerToModPolyVec(t)

	t_1_packed := dil.bitPackPolyVec(t_1, 10)
	shake := append(rho, t_1_packed...)
	tr := dil.shake256(shake, 32)

	pk = append(rho, t_1_packed...)

	sk = append(rho, K...)
	sk = append(sk, tr...)
	sk = append(sk, dil.bitPackAlteredPolyVec(s_1, dil.eta, dil.sBits)...)
	sk = append(sk, dil.bitPackAlteredPolyVec(s_2, dil.eta, dil.sBits)...)
	sk = append(sk, dil.bitPackAlteredPolyVec(t_0, 1<<d, d)...)
	return
}

func (dil *Dilithium) Sign(sk []byte, message []byte) (sigma []byte) {
	rho := sk[:32]
	K := sk[32:64]
	tr := sk[64:96]
	s_1 := dil.modPMPolyVec(
		dil.bitUnpackAlteredPolyVec(
			sk[96:dil.s1Bytes+96],
			dil.eta,
			dil.sBits),
		q)
	s_2 := dil.modPMPolyVec(
		dil.bitUnpackAlteredPolyVec(
			sk[96+dil.s1Bytes:dil.s2Bytes+dil.s1Bytes+96],
			dil.eta,
			dil.sBits),
		q)
	t_0 := dil.modPMPolyVec(
		dil.bitUnpackAlteredPolyVec(
			sk[96+dil.s2Bytes+dil.s1Bytes:],
			1<<d,
			d),
		1<<d)

	A_hat := dil.expandA(rho)
	s_1_hat := dil.nttPolyVec(s_1)
	s_2_hat := dil.nttPolyVec(s_2)
	t_0_hat := dil.nttPolyVec(t_0)

	shake := append(myio.Copy(tr), message...)

	mi := dil.shake256(shake, 64)
	shake = append(myio.Copy(K), mi...)
	rho_dash := dil.shake256(shake, 64)
	kappa := -dil.l

	for {
		kappa += dil.l

		y := dil.expandMask(rho_dash, kappa)
		y_hat := dil.nttPolyVec(y)

		w := make([][]int, dil.k)
		for i := 0; i < dil.k; i++ {
			w[i] = dil.invNtt(dil.mulPolyVec(A_hat[i], y_hat))
		}
		w_1 := dil.highBitsPolyVec(w, 2*dil.gammaTwo)

		w_1_packed := dil.bitPackPolyVec(w_1, dil.w_1Bits)

		shake = append(mi, w_1_packed...)
		c_wave := dil.shake256(shake, 32)

		c := dil.sampleInBall(c_wave)
		c_hat := dil.ntt(c)

		cs_1 := dil.invNttPolyVec(dil.scalePolyVecByPoly(s_1_hat, c_hat))
		cs_2 := dil.invNttPolyVec(dil.scalePolyVecByPoly(s_2_hat, c_hat))

		z := dil.modPMPolyVec(dil.addPolyVec(y, cs_1), q)

		w_sub_cs_2 := dil.modPMPolyVec(dil.subPolyVec(w, cs_2), q)

		r_0 := dil.lowBitsPolyVec(dil.modPMPolyVec(w_sub_cs_2, q), 2*dil.gammaTwo)

		if dil.checkNormPolyVec(z, dil.gammaOne-dil.beta) || dil.checkNormPolyVec(r_0, dil.gammaTwo-dil.beta) {
			continue
		}

		ct_0 := dil.modPMPolyVec(dil.invNttPolyVec(dil.scalePolyVecByPoly(t_0_hat, c_hat)), q)
		ct_0_inv := dil.modPMPolyVec(dil.inversePolyVec(ct_0), q)
		w_sub_cs_2_add_ct_0 := dil.modPMPolyVec(dil.addPolyVec(w_sub_cs_2, ct_0), q)

		h := dil.makeHintPolyVec(ct_0_inv, w_sub_cs_2_add_ct_0, 2*dil.gammaTwo)

		ones := 0
		for i := 0; i < len(h); i++ {
			for j := 0; j < n; j++ {
				if h[i][j] == 1 {
					ones++
				}
			}
		}
		if dil.checkNormPolyVec(dil.modPMPolyVec(ct_0, q), dil.gammaTwo) {
			continue
		}
		if ones > dil.omega {
			continue
		}

		z_packed := dil.bitPackAlteredPolyVec(z, dil.gammaOne, dil.zBits)
		h_packed := dil.bitPackHint(h)
		sigma = append(c_wave, z_packed...)
		sigma = append(sigma, h_packed...)
		break
	}
	return
}

func (dil *Dilithium) Verify(pk, message, sigma []byte) (verified bool) {
	rho := pk[:32]
	t_1_bytes := pk[32:]
	c_wave := sigma[:32]
	z := dil.bitUnpackAlteredPolyVec(sigma[32:dil.zBytes+32], dil.gammaOne, dil.zBits)
	h := dil.bitUnpackHint(sigma[dil.zBytes+32:])
	t_1 := dil.bitUnpackPolyVec(t_1_bytes, 10)

	A_hat := dil.expandA(rho)

	shake := append(rho, t_1_bytes...)
	shake = append(dil.shake256(shake, 32), message...)
	mi := dil.shake256(shake, 64)

	c := dil.sampleInBall(c_wave)
	c_hat := dil.ntt(c)

	z_hat := dil.nttPolyVec(z)

	Az := make([][]int, dil.k)
	for i := 0; i < dil.k; i++ {
		Az[i] = dil.mulPolyVec(A_hat[i], z_hat)
	}

	t_1_2_to_d := dil.modPMPolyVec(dil.scalePolyVecByInt(t_1, 1<<d), q)
	t_1_2_to_d_hat := dil.nttPolyVec(t_1_2_to_d)
	ct_1 := dil.scalePolyVecByPoly(t_1_2_to_d_hat, c_hat)

	r := dil.invNttPolyVec(dil.subPolyVec(Az, ct_1))

	w_1 := dil.useHintPolyVec(h, dil.modPMPolyVec(r, q), 2*dil.gammaTwo)
	shake = append(mi, dil.bitPackPolyVec(w_1, dil.w_1Bits)...)
	verified = dil.BytesEqual(c_wave, dil.shake256(shake, 32))
	return
}
