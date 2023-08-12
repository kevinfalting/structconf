package confhandler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/kevinfalting/structconf/stronf"
)

// RSA is a struct that holds a pair of RSA keys and provides methods for
// handling RSA cryptographic operations. This handler is best used as the final
// handler.
type RSAHandler struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey

	// Label is an optional byte slice used during encryption and decryption
	// with RSA, and enhances message integrity checks. The same Label must be
	// used for both encryption and decryption if one is set.
	// Note: The Label is not secret and does not get encrypted.
	Label []byte
}

var _ stronf.Handler = (*RSAHandler)(nil)

// Handle applies the RSA decryption process to a given field value. It checks
// for the "secret" tag and decrypts the value if found. If a previous handler
// has attempted to set a value, represented by the interimValue, this handler
// will decrypt the interimValue, otherwise it will decrypt the field value.
func (r *RSAHandler) Handle(ctx context.Context, f stronf.Field, interimValue any) (any, error) {
	_, ok := f.LookupTag("conf", "secret")
	if !ok {
		return nil, nil
	}

	val := f.Value()
	if interimValue != nil {
		val = interimValue
	}

	ciphertext, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("field %q or interimValue must be of type string, got: %T", f.Name(), val)
	}

	if len(ciphertext) == 0 {
		// If there is nothing to decrypt, that's okay too.
		return nil, nil
	}

	plaintext, err := r.Decrypt([]byte(ciphertext))
	if err != nil {
		return nil, err
	}

	return string(plaintext), nil
}

// Decrypt decodes the provided ciphertext from base64 and decrypts it using RSA OAEP.
func (r *RSAHandler) Decrypt(cipherText []byte) ([]byte, error) {
	encryptedBytes := make([]byte, base64.StdEncoding.DecodedLen(len(cipherText)))
	n, err := base64.StdEncoding.Decode(encryptedBytes, cipherText)
	if err != nil {
		return nil, err
	}

	encryptedBytes = encryptedBytes[:n]

	decryptedBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		nil,
		r.PrivateKey,
		encryptedBytes,
		r.Label,
	)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}

// Encrypt encrypts the plaintext using RSA OAEP, then encodes the ciphertext into base64.
func (r *RSAHandler) Encrypt(plaintext []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		r.PublicKey,
		plaintext,
		r.Label,
	)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedBytes)))
	base64.StdEncoding.Encode(cipherText, encryptedBytes)

	return cipherText, nil
}

// NewPEMKeyPair generates a new pair of RSA keys with the provided bit size,
// sets them as the PrivateKey and PublicKey of the RSAHandler, and returns them
// in PEM format. A typical safe bit size is 2048.
func (r *RSAHandler) NewPEMKeyPair(bitSize int) (privateKey, publicKey []byte, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	r.PrivateKey = privKey

	pubKey := privKey.PublicKey
	r.PublicKey = &pubKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

	return privateKeyPEM, publicKeyPEM, nil
}

// PublicKeyFromPEM takes a PEM-formatted RSA public key as input, decodes it,
// sets it as the RSAHandler's public key, and returns it. If the provided key
// is not a valid RSA public key, it returns an error.
func (r *RSAHandler) PublicKeyFromPEM(pubkeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubkeyPEM)
	if block == nil {
		return nil, errors.New("failed to pem.Decode public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	r.PublicKey = publicKey
	return publicKey, nil
}

// PrivateKeyFromPEM takes a PEM-formatted RSA private key as input, decodes it,
// sets it as the RSAHandler's private key, and returns it. If the provided key
// is not a valid RSA private key, it returns an error.
func (r *RSAHandler) PrivateKeyFromPEM(privkeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privkeyPEM)
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	r.PrivateKey = priv
	return priv, nil
}
