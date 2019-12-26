package utility

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	_log "github.com/stevenb256/log"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

// NonceSize size of nonce
const NonceSize = 24

// KeySize lengtho of crypto key
const KeySize = 32

// ErrInvalidCryptoKey invalid crypto key
var ErrInvalidCryptoKey = _log.NewError(100, "crypto", "invalid crypto key length")

// ErrCantOpenSealedBytes can't open sealed bytes; crypto problem
var ErrCantOpenSealedBytes = _log.NewError(101, "crypto", "unable to open/unseal bytes")

// ErrCantDecryptBytes can't decrypt bytes encrypted in other function
var ErrCantDecryptBytes = _log.NewError(102, "crypto", "unable to decrypt bytes")

// GenerateCryptoKeys returns public, private keys or error
func GenerateCryptoKeys() (*[KeySize]byte, *[KeySize]byte, error) {
	return box.GenerateKey(rand.Reader)
}

// CryptoKeyFromBase64 get crypto key from base64 string
func CryptoKeyFromBase64(key64 string) (*[KeySize]byte, error) {
	buf, err := base64.StdEncoding.DecodeString(key64)
	if _log.Check(err) {
		return nil, err
	}
	if len(buf) != KeySize {
		return nil, _log.Fail(ErrInvalidCryptoKey, key64)
	}
	var key [KeySize]byte
	copy((key)[:], buf)
	return &key, nil
}

// SealBytes encrypts/signs buffer with a public key of recipient and private key
// of the sender
func SealBytes(buf []byte, public, private *[KeySize]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	io.ReadFull(rand.Reader, nonce[:])
	return box.Seal(nonce[:], buf, &nonce, public, private), nil
}

// OpenSealedBytes - decrypts bytes with public key of the sender and
// private key of the recipient
func OpenSealedBytes(buf []byte, public, private *[KeySize]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	_log.Assert(len(buf) >= len(nonce))
	copy(nonce[:], buf[:NonceSize])
	clear, b := box.Open(nil, buf[NonceSize:], &nonce, public, private)
	if false == b {
		return nil, _log.Fail(ErrCantOpenSealedBytes)
	}
	return clear, nil
}

// EncryptBytes used to just encrypt bytes with a random key
func EncryptBytes(in []byte, key *[KeySize]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	io.ReadFull(rand.Reader, nonce[:])
	out := make([]byte, NonceSize)
	copy(out, nonce[:])
	return secretbox.Seal(out, in, &nonce, key), nil
}

// DecryptBytes used to decrypt bytes with a key used in EncryptBytes
func DecryptBytes(in []byte, key *[KeySize]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	_log.Assert(len(in) >= NonceSize)
	copy(nonce[:], in[:NonceSize])
	out, worked := secretbox.Open(nil, in[NonceSize:], &nonce, key)
	if false == worked {
		return nil, _log.Fail(ErrCantDecryptBytes)
	}
	return out, nil
}
