package dilithium

const (
	Q        = 8380417
	N        = 256
	Tau      = 39
	K        = 4
	L        = 4
	Eta      = 2
	D        = 13
	SBytes   = (N * 3 / 8) * K
	ZBytes   = (N * 18 / 8) * K
	GammaOne = 1 << 17
	GammaTwo = (Q - 1) / 88
	Omega    = 80
	Beta     = Tau * Eta
)
