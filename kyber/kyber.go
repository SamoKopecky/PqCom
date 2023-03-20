package kyber

import (
	"github.com/SamoKopecky/pqcom/main/io"
)

const (
	q             = 3329
	n             = 256
	sharedKeySize = 32
)

type Kyber struct {
	k      int
	eta2   int
	eta1   int
	du     int
	dv     int
	PkSize int
	SkSize int
	CtSize int
}

func Kyber512() Kyber {
	return Kyber{
		k:      2,
		eta2:   2,
		eta1:   3,
		du:     10,
		dv:     4,
		PkSize: 800,
		SkSize: 1632,
		CtSize: 768,
	}
}

func Kyber768() Kyber {
	return Kyber{
		k:      3,
		eta2:   2,
		eta1:   2,
		du:     10,
		dv:     4,
		PkSize: 1184,
		SkSize: 2400,
		CtSize: 1088,
	}
}

func Kyber1024() Kyber {
	return Kyber{
		k:      4,
		eta2:   2,
		eta1:   2,
		du:     11,
		dv:     5,
		PkSize: 1568,
		SkSize: 3168,
		CtSize: 1568,
	}
}

func (kyb *Kyber) cpapkeKeyGen() (pk []byte, sk []byte) {
	d := kyb.randBytes(32)
	t_hat := make([][]int, kyb.k)
	localN := byte(0)
	rho, sigma := hash64(d)

	A_hat := kyb.genPolyMat(rho, false)
	s_hat := kyb.randPolyVec(sigma, &localN, kyb.eta1)
	e_hat := kyb.randPolyVec(sigma, &localN, kyb.eta1)

	kyb.nttPolyVec(s_hat)
	kyb.nttPolyVec(e_hat)

	for i := 0; i < kyb.k; i++ {
		t_hat[i] = kyb.mulPolyVec(A_hat[i], s_hat)
	}
	t_hat = kyb.addPolyVec(t_hat, e_hat)

	kyb.modPlusPolyVec(s_hat)
	kyb.modPlusPolyVec(t_hat)

	sk = kyb.encodePolyVec(s_hat, 12)
	pk = kyb.encodePolyVec(t_hat, 12)
	pk = append(pk, rho...)
	return
}

func (kyb *Kyber) cpapkeEnc(pk []byte, m []byte, randomCoins []byte) (c []byte) {
	c1 := []byte{}
	localN := byte(0)

	t_hat := kyb.decodePolyVec(pk, 12)
	rho := pk[len(pk)-32:]

	A_hat := kyb.genPolyMat(rho, true)
	r_hat := kyb.randPolyVec(randomCoins, &localN, kyb.eta1)
	e1 := kyb.randPolyVec(randomCoins, &localN, kyb.eta2)
	e2 := kyb.randPoly(randomCoins, &localN, kyb.eta2)

	kyb.nttPolyVec(r_hat)

	u := make([][]int, kyb.k)
	for i := 0; i < kyb.k; i++ {
		u[i] = kyb.mulPolyVec(A_hat[i], r_hat)
	}

	kyb.invNttPolyVec(u)
	u = kyb.addPolyVec(u, e1)

	parsed_m := kyb.decompress(kyb.decode(m, 1), 1)
	v := kyb.mulPolyVec(t_hat, r_hat)
	kyb.invNtt(v)
	v = kyb.polyAdd(v, e2)
	v = kyb.polyAdd(v, parsed_m)

	for i := 0; i < kyb.k; i++ {
		c1 = append(c1, kyb.encode(kyb.compress(u[i], kyb.du), kyb.du)...)
	}

	c2 := kyb.encode(kyb.compress(v, kyb.dv), kyb.dv)
	c = append(c1, c2...)

	return
}

func (kyb *Kyber) cpapkeDec(sk []byte, c []byte) (m []byte) {
	u_hat := make([][]int, kyb.k)

	c2 := c[kyb.du*kyb.k*n/8:]
	u_decoded := kyb.decodePolyVec(c, kyb.du)

	for i := 0; i < kyb.k; i++ {
		u_hat[i] = kyb.decompress(u_decoded[i], kyb.du)
	}

	kyb.nttPolyVec(u_hat)

	v := kyb.decompress(kyb.decode(c2, kyb.dv), kyb.dv)
	s_hat := kyb.decodePolyVec(sk, 12)

	s_hat_u_hat := kyb.mulPolyVec(s_hat, u_hat)
	kyb.invNtt(s_hat_u_hat)
	first_m := kyb.polySub(v, s_hat_u_hat)

	m = kyb.encode(kyb.compress(first_m, 1), 1)
	return
}

func (kyb *Kyber) CcakemKeyGen() (pk, sk []byte) {
	z := kyb.randBytes(32)

	pk, sk_dash := kyb.cpapkeKeyGen()

	sk = []byte{}
	sk = append(sk, sk_dash...)
	sk = append(sk, pk...)
	sk = append(sk, hash32(pk)...)
	sk = append(sk, z...)
	return
}

func (kyb *Kyber) CcakemEnc(pk []byte) (c, key []byte) {
	m := hash32(kyb.randBytes(32))

	g_input := []byte{}
	g_input = append(g_input, m...)
	g_input = append(g_input, hash32(pk)...)

	K_dash, r := hash64(g_input)
	c = kyb.cpapkeEnc(pk, m, r)

	kdf_input := []byte{}
	kdf_input = append(kdf_input, K_dash...)
	kdf_input = append(kdf_input, hash32(c)...)
	key = kdf(kdf_input, 32)
	return
}

func (kyb *Kyber) CcakemDec(c, sk []byte) (key []byte) {
	keySize := 12 * kyb.k * n / 8
	skCopy := io.Copy(sk)
	pk := skCopy[keySize : keySize*2+32]
	hash := skCopy[keySize*2+32 : keySize*2+64]
	z := skCopy[keySize*2+64:]

	m_dash := kyb.cpapkeDec(sk, c)

	g_input := []byte{}
	g_input = append(g_input, m_dash...)
	g_input = append(g_input, hash...)
	k_dash, r_dash := hash64(g_input)

	c_dash := kyb.cpapkeEnc(pk, m_dash, r_dash)
	hash_c := hash32(c)

	kdf_input := []byte{}
	if kyb.BytesEqual(c, c_dash) {
		kdf_input = append(kdf_input, k_dash...)
		kdf_input = append(kdf_input, hash_c...)
		key = kdf(kdf_input, sharedKeySize)
	} else {
		kdf_input = append(kdf_input, z...)
		kdf_input = append(kdf_input, hash_c...)
		key = kdf(kdf_input, sharedKeySize)
	}
	return
}
