package envelope

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return key
}

func TestKeyring_EncryptDecrypt_RoundTrip(t *testing.T) {
	k, err := NewKeyring(1, map[uint8][]byte{1: newKey(t)})
	require.NoError(t, err)

	blob, err := k.Encrypt([]byte("alice@example.com"))
	require.NoError(t, err)

	plaintext, err := k.Decrypt(blob)
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", string(plaintext))
}

func TestKeyring_EncryptUsesCurrentVersion(t *testing.T) {
	k, err := NewKeyring(2, map[uint8][]byte{1: newKey(t), 2: newKey(t)})
	require.NoError(t, err)

	blob, err := k.Encrypt([]byte("alice@example.com"))
	require.NoError(t, err)

	version, err := VersionOf(blob)
	require.NoError(t, err)
	assert.Equal(t, uint8(2), version)
	assert.Equal(t, uint8(2), k.Current())
}

func TestKeyring_DecryptOldVersion(t *testing.T) {
	key1 := newKey(t)
	v1, err := NewKeyring(1, map[uint8][]byte{1: key1})
	require.NoError(t, err)
	blob, err := v1.Encrypt([]byte("010-1234-5678"))
	require.NoError(t, err)

	v2, err := NewKeyring(2, map[uint8][]byte{1: key1, 2: newKey(t)})
	require.NoError(t, err)

	plaintext, err := v2.Decrypt(blob)
	require.NoError(t, err)
	assert.Equal(t, "010-1234-5678", string(plaintext))
}

func TestKeyring_EncryptNonceUnique(t *testing.T) {
	k, err := NewKeyring(1, map[uint8][]byte{1: newKey(t)})
	require.NoError(t, err)

	a, err := k.Encrypt([]byte("same"))
	require.NoError(t, err)
	b, err := k.Encrypt([]byte("same"))
	require.NoError(t, err)

	assert.NotEqual(t, a, b)
}

func TestKeyring_DecryptUnknownVersion(t *testing.T) {
	v2, err := NewKeyring(2, map[uint8][]byte{2: newKey(t)})
	require.NoError(t, err)
	blob, err := v2.Encrypt([]byte("alice@example.com"))
	require.NoError(t, err)

	v1, err := NewKeyring(1, map[uint8][]byte{1: newKey(t)})
	require.NoError(t, err)

	_, err = v1.Decrypt(blob)
	assert.ErrorIs(t, err, ErrUnknownKeyVersion)
}

func TestKeyring_DecryptTooShort(t *testing.T) {
	k, err := NewKeyring(1, map[uint8][]byte{1: newKey(t)})
	require.NoError(t, err)

	short := make([]byte, 1+12+15)
	short[0] = 1

	_, err = k.Decrypt(short)
	assert.ErrorIs(t, err, ErrCiphertextTooShort)

	_, err = k.Decrypt(nil)
	assert.ErrorIs(t, err, ErrCiphertextTooShort)
}

func TestKeyring_DecryptTampered(t *testing.T) {
	k, err := NewKeyring(1, map[uint8][]byte{1: newKey(t)})
	require.NoError(t, err)
	blob, err := k.Encrypt([]byte("alice@example.com"))
	require.NoError(t, err)

	blob[len(blob)-1] ^= 0xff

	_, err = k.Decrypt(blob)
	assert.Error(t, err)
}

func TestNewKeyring_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		current uint8
		keys    map[uint8][]byte
	}{
		{name: "31바이트 키", current: 1, keys: map[uint8][]byte{1: make([]byte, 31)}},
		{name: "current 버전 키 없음", current: 2, keys: map[uint8][]byte{1: newKey(t)}},
		{name: "current가 0", current: 0, keys: map[uint8][]byte{0: newKey(t)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKeyring(tt.current, tt.keys)
			assert.ErrorIs(t, err, ErrInvalidKey)
		})
	}
}
