// Package envelope 는 Vault Transit KEK 로 보호하는 DEK 를 써서 데이터를 암호화하는
// envelope encryption 예제다.
package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	keySize   = 32 // AES-256
	nonceSize = 12 // GCM 표준 nonce 크기
	tagSize   = 16 // GCM 인증 tag 크기
)

var (
	// ErrUnknownKeyVersion 은 keyring 에 없는 버전으로 만든 암호문을 복호화하려 할 때 반환된다.
	ErrUnknownKeyVersion = errors.New("unknown key version")
	// ErrCiphertextTooShort 는 암호문이 버전·nonce·tag 를 담기에도 짧을 때 반환된다.
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	// ErrInvalidKey 는 키 길이나 current 버전이 잘못됐을 때 반환된다.
	ErrInvalidKey = errors.New("invalid key")
)

// Keyring 은 버전별 DEK 를 들고 있다. 암호화는 항상 current 버전으로 하고,
// 복호화는 암호문 첫 바이트에 적힌 버전의 키로 한다.
//
// 암호문 포맷: [1바이트 버전][12바이트 nonce][AES-GCM 암호문 + 16바이트 tag]
type Keyring struct {
	current uint8
	aeads   map[uint8]cipher.AEAD
}

// NewKeyring 은 버전별 32바이트 키로 Keyring 을 만든다. current 는 1 이상이고 keys 에 있어야 한다.
func NewKeyring(current uint8, keys map[uint8][]byte) (*Keyring, error) {
	if current == 0 {
		return nil, fmt.Errorf("%w: current version must be 1 or greater", ErrInvalidKey)
	}
	if _, ok := keys[current]; !ok {
		return nil, fmt.Errorf("%w: no key for current version %d", ErrInvalidKey, current)
	}

	aeads := make(map[uint8]cipher.AEAD, len(keys))
	for version, key := range keys {
		if len(key) != keySize {
			return nil, fmt.Errorf("%w: version %d key is %d bytes, want %d", ErrInvalidKey, version, len(key), keySize)
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		aead, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		aeads[version] = aead
	}

	return &Keyring{current: current, aeads: aeads}, nil
}

// Current 는 암호화에 쓰는 DEK 버전을 돌려준다.
func (k *Keyring) Current() uint8 {
	return k.current
}

// Encrypt 는 current 버전 키로 평문을 암호화한다. 호출할 때마다 새 nonce 를 쓴다.
func (k *Keyring) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	blob := make([]byte, 0, 1+nonceSize+len(plaintext)+tagSize)
	blob = append(blob, k.current)
	blob = append(blob, nonce...)
	return k.aeads[k.current].Seal(blob, nonce, plaintext, nil), nil
}

// Decrypt 는 암호문 첫 바이트의 버전 키로 복호화하고, 변조됐으면 에러를 돌려준다.
func (k *Keyring) Decrypt(blob []byte) ([]byte, error) {
	version, err := VersionOf(blob)
	if err != nil {
		return nil, err
	}
	aead, ok := k.aeads[version]
	if !ok {
		return nil, fmt.Errorf("%w: %d", ErrUnknownKeyVersion, version)
	}

	nonce := blob[1 : 1+nonceSize]
	return aead.Open(nil, nonce, blob[1+nonceSize:], nil)
}

// VersionOf 는 암호문을 만든 DEK 버전을 돌려준다.
func VersionOf(blob []byte) (uint8, error) {
	if len(blob) < 1+nonceSize+tagSize {
		return 0, fmt.Errorf("%w: %d bytes", ErrCiphertextTooShort, len(blob))
	}
	return blob[0], nil
}
