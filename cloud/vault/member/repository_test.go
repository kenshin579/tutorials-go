package member

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/cloud/vault/envelope"
)

var alice = Member{ID: "alice", Email: "alice@example.com", Phone: "010-1234-5678"}

func newKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return key
}

func newKeyring(t *testing.T, current uint8, keys map[uint8][]byte) *envelope.Keyring {
	t.Helper()
	k, err := envelope.NewKeyring(current, keys)
	require.NoError(t, err)
	return k
}

// fieldVersions 는 파일에 저장된 모든 암호문 필드의 DEK 버전을 모은다.
func fieldVersions(t *testing.T, path string) []uint8 {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var f struct {
		Members []map[string]string `json:"members"`
	}
	require.NoError(t, json.Unmarshal(raw, &f))

	var versions []uint8
	for _, m := range f.Members {
		for _, field := range []string{"email", "phone"} {
			blob, err := base64.StdEncoding.DecodeString(m[field])
			require.NoError(t, err)
			v, err := envelope.VersionOf(blob)
			require.NoError(t, err)
			versions = append(versions, v)
		}
	}
	return versions
}

func TestRepository_PutGet_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "members.json")
	repo := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))

	require.NoError(t, repo.Put(alice))

	got, err := repo.Get("alice")
	require.NoError(t, err)
	assert.Equal(t, alice, got)
}

func TestRepository_FileHasNoPlaintext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	repo := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))
	require.NoError(t, repo.Put(alice))

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), alice.Email)
	assert.NotContains(t, string(raw), alice.Phone)
}

func TestRepository_GetNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	repo := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))
	require.NoError(t, repo.Put(alice))

	_, err := repo.Get("bob")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepository_GetMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	repo := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))

	_, err := repo.Get("alice")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepository_Reencrypt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	key1 := newKey(t)
	v1 := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: key1}))
	bob := Member{ID: "bob", Email: "bob@example.com", Phone: "010-9876-5432"}
	require.NoError(t, v1.Put(alice))
	require.NoError(t, v1.Put(bob))

	v2 := NewRepository(path, newKeyring(t, 2, map[uint8][]byte{1: key1, 2: newKey(t)}))
	n, err := v2.Reencrypt()
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []uint8{2, 2, 2, 2}, fieldVersions(t, path))

	n, err = v2.Reencrypt()
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	got, err := v2.Get("bob")
	require.NoError(t, err)
	assert.Equal(t, bob, got)
}

func TestRepository_ReencryptEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	repo := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))

	n, err := repo.Reencrypt()
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestRepository_GetWithWrongKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.json")
	require.NoError(t, NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)})).Put(alice))

	// 같은 버전 번호(1)지만 다른 키: Vault 를 재시작해 KEK·DEK 를 새로 만든 상황
	other := NewRepository(path, newKeyring(t, 1, map[uint8][]byte{1: newKey(t)}))
	got, err := other.Get("alice")
	assert.Error(t, err)
	assert.Equal(t, Member{}, got)
}
