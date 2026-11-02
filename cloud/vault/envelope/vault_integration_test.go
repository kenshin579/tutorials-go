package envelope

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcvault "github.com/testcontainers/testcontainers-go/modules/vault"
)

const vaultImage = "hashicorp/vault:2.2.1"

// startVault 는 Transit KEK 와 policy 를 준비한 dev Vault 를 띄우고 root client 를 돌려준다.
func startVault(t *testing.T) *vaultapi.Client {
	t.Helper()
	if testing.Short() {
		t.Skip("Docker 가 필요한 통합 테스트")
	}
	ctx := context.Background()

	ctr, err := tcvault.Run(ctx, vaultImage,
		tcvault.WithToken("root"),
		tcvault.WithInitCommand("secrets enable transit", "write -f transit/keys/member-api-kek"),
	)
	t.Cleanup(func() {
		if ctr != nil {
			_ = ctr.Terminate(context.Background())
		}
	})
	require.NoError(t, err)

	addr, err := ctr.HttpHostAddress(ctx)
	require.NoError(t, err)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	root, err := vaultapi.NewClient(cfg)
	require.NoError(t, err)
	root.SetToken("root")

	for _, name := range []string{"member-api", "member-api-migration"} {
		rules, err := os.ReadFile("../policies/" + name + ".hcl")
		require.NoError(t, err)
		require.NoError(t, root.Sys().PutPolicy(name, string(rules)))
	}
	return root
}

// newTokenClient 는 policy 하나만 가진 token 으로 동작하는 client 를 만든다.
func newTokenClient(t *testing.T, root *vaultapi.Client, policy string) *vaultapi.Client {
	t.Helper()
	secret, err := root.Auth().Token().Create(&vaultapi.TokenCreateRequest{
		Policies: []string{policy},
		TTL:      "1h",
	})
	require.NoError(t, err)

	c, err := root.Clone()
	require.NoError(t, err)
	c.SetToken(secret.Auth.ClientToken)
	return c
}

func readDEKDoc(t *testing.T, root *vaultapi.Client) map[string]interface{} {
	t.Helper()
	secret, err := root.KVv2("secret").Get(context.Background(), "member-api/dek")
	require.NoError(t, err)
	return secret.Data
}

func assertStatus(t *testing.T, err error, want int) {
	t.Helper()
	var respErr *vaultapi.ResponseError
	require.True(t, errors.As(err, &respErr), "want *vaultapi.ResponseError, got %v", err)
	assert.Equal(t, want, respErr.StatusCode)
}

func TestInitDEK_Idempotent(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")

	created, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	assert.True(t, created)

	created, err = InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	assert.False(t, created)

	assert.Equal(t, "1", readDEKDoc(t, root)["current_version"])
}

func TestRotateDEK_AddsVersion(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	wrappedV1 := readDEKDoc(t, root)["wrapped_v1"]

	version, err := RotateDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	assert.Equal(t, uint8(2), version)

	doc := readDEKDoc(t, root)
	assert.Equal(t, "2", doc["current_version"])
	assert.Equal(t, wrappedV1, doc["wrapped_v1"], "기존 버전은 덮어쓰지 않는다")
	assert.NotEmpty(t, doc["wrapped_v2"])
}

func TestLoadKeyring_DecryptsAcrossRotation(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")
	app := newTokenClient(t, root, "member-api")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	v1, err := LoadKeyring(ctx, app, DefaultConfig)
	require.NoError(t, err)
	blob, err := v1.Encrypt([]byte("alice@example.com"))
	require.NoError(t, err)

	_, err = RotateDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)

	v2, err := LoadKeyring(ctx, app, DefaultConfig)
	require.NoError(t, err)
	assert.Equal(t, uint8(2), v2.Current())

	plaintext, err := v2.Decrypt(blob)
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", string(plaintext))

	fresh, err := v2.Encrypt([]byte("bob@example.com"))
	require.NoError(t, err)
	version, err := VersionOf(fresh)
	require.NoError(t, err)
	assert.Equal(t, uint8(2), version)
}

func TestLoadKeyring_NotInitialized(t *testing.T) {
	root := startVault(t)
	app := newTokenClient(t, root, "member-api")

	_, err := LoadKeyring(context.Background(), app, DefaultConfig)
	assert.ErrorIs(t, err, ErrDEKNotInitialized)
}

func TestWriteDEK_StaleCASFails(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)

	err = writeDEK(ctx, migration, DefaultConfig, map[string]interface{}{
		"current_version": "9",
		"wrapped_v9":      "vault:v1:stale",
	}, 0)
	require.Error(t, err)

	assert.Equal(t, "1", readDEKDoc(t, root)["current_version"])
}

func TestPolicy_AppTokenCannotRotate(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")
	app := newTokenClient(t, root, "member-api")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)

	_, err = RotateDEK(ctx, app, DefaultConfig)
	assertStatus(t, err, http.StatusForbidden)
}

func TestPolicy_MigrationTokenCannotRotateKEK(t *testing.T) {
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")

	_, err := migration.Logical().WriteWithContext(context.Background(), "transit/keys/member-api-kek/rotate", nil)
	assertStatus(t, err, http.StatusForbidden)
}

func TestRotateDEK_DoesNotOverwriteAfterRollback(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	_, err = RotateDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)

	// 운영자가 current_version 을 1 로 되돌린 상황
	doc := readDEKDoc(t, root)
	wrappedV2 := doc["wrapped_v2"]
	doc["current_version"] = "1"
	_, err = root.KVv2("secret").Put(ctx, "member-api/dek", doc)
	require.NoError(t, err)

	version, err := RotateDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	assert.Equal(t, uint8(3), version)

	doc = readDEKDoc(t, root)
	assert.Equal(t, wrappedV2, doc["wrapped_v2"], "기존 v2 를 덮어쓰지 않는다")
	assert.Equal(t, "3", doc["current_version"])
}

func TestInitDEK_DeletedDocument(t *testing.T) {
	ctx := context.Background()
	root := startVault(t)
	migration := newTokenClient(t, root, "member-api-migration")
	app := newTokenClient(t, root, "member-api")

	_, err := InitDEK(ctx, migration, DefaultConfig)
	require.NoError(t, err)
	require.NoError(t, root.KVv2("secret").Delete(ctx, "member-api/dek"))

	_, err = InitDEK(ctx, migration, DefaultConfig)
	assert.ErrorIs(t, err, ErrDEKDeleted)

	_, err = LoadKeyring(ctx, app, DefaultConfig)
	assert.ErrorIs(t, err, ErrDEKDeleted)
}
