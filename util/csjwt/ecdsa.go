package csjwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	// Sadly this is missing from crypto/ecdsa compared to crypto/rsa
	ErrECDSAVerification = errors.New("crypto/ecdsa: verification error")
)

// Implements the ECDSA family of signing methods signing methods
type SigningMethodECDSA struct {
	Name      string
	Hash      crypto.Hash
	KeySize   int
	CurveBits int
}

// Specific instances for EC256 and company
var (
	SigningMethodES256 *SigningMethodECDSA
	SigningMethodES384 *SigningMethodECDSA
	SigningMethodES512 *SigningMethodECDSA
)

func init() {
	// ES256
	SigningMethodES256 = &SigningMethodECDSA{"ES256", crypto.SHA256, 32, 256}
	RegisterSigningMethod(SigningMethodES256.Alg(), func() SigningMethod {
		return SigningMethodES256
	})

	// ES384
	SigningMethodES384 = &SigningMethodECDSA{"ES384", crypto.SHA384, 48, 384}
	RegisterSigningMethod(SigningMethodES384.Alg(), func() SigningMethod {
		return SigningMethodES384
	})

	// ES512
	SigningMethodES512 = &SigningMethodECDSA{"ES512", crypto.SHA512, 66, 521}
	RegisterSigningMethod(SigningMethodES512.Alg(), func() SigningMethod {
		return SigningMethodES512
	})
}

func (m *SigningMethodECDSA) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this verify method, key must be an ecdsa.PublicKey struct
func (m *SigningMethodECDSA) Verify(signingString, signature []byte, key interface{}) error {
	// Decode the signature
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Get the key
	var ecdsaKey *ecdsa.PublicKey
	switch k := key.(type) {
	case *ecdsa.PublicKey:
		ecdsaKey = k
	default:
		return ErrInvalidKey
	}

	if len(sig) != 2*m.KeySize {
		return ErrECDSAVerification
	}

	r := big.NewInt(0).SetBytes(sig[:m.KeySize])
	s := big.NewInt(0).SetBytes(sig[m.KeySize:])

	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return err
	}

	// Verify the signature
	err = ErrECDSAVerification
	if ecdsa.Verify(ecdsaKey, hasher.Sum(nil), r, s) {
		err = nil
	}
	return err
}

// Implements the Sign method from SigningMethod
// For this signing method, key must be an ecdsa.PrivateKey struct
func (m *SigningMethodECDSA) Sign(signingString []byte, key interface{}) ([]byte, error) {
	// Get the key
	var ecdsaKey *ecdsa.PrivateKey
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		ecdsaKey = k
	default:
		return nil, ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return nil, ErrHashUnavailable
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return nil, err
	}

	// Sign the string and return r, s
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}

	curveBits := ecdsaKey.Curve.Params().BitSize

	if m.CurveBits != curveBits {
		return nil, ErrInvalidKey
	}

	keyBytes := curveBits / 8
	if curveBits%8 > 0 {
		keyBytes += 1
	}

	// We serialize the outpus (r and s) into big-endian byte arrays and pad
	// them with zeros on the left to make sure the sizes work out. Both arrays
	// must be keyBytes long, and the output must be 2*keyBytes long.
	rBytes := r.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := s.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)

	return EncodeSegment(out), nil

}