package envelope

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// KV v2 문서에 DEK 를 저장하는 필드 이름이다.
//
//	current_version = "2"          // 암호화에 쓰는 버전
//	wrapped_v1      = "vault:v1:…" // KEK 로 wrap 된 DEK (버전별로 남긴다)
//	wrapped_v2      = "vault:v1:…"
const (
	currentVersionField = "current_version"
	wrappedFieldPrefix  = "wrapped_v"
	dekBits             = 256
)

var (
	// ErrDEKNotInitialized 는 KV 에 DEK 문서가 아직 없을 때 LoadKeyring 이 반환한다.
	ErrDEKNotInitialized = errors.New("DEK not initialized")
	// ErrDEKDeleted 는 DEK 문서의 최신 버전이 삭제(soft delete)된 상태일 때 반환된다.
	// 옛 데이터가 그 DEK 에 의존할 수 있으므로 새로 만들지 않고 운영자에게 판단을 넘긴다.
	ErrDEKDeleted = errors.New("DEK document is deleted")
)

// Config 는 KEK 이름과 wrap 된 DEK 를 저장할 KV v2 위치다.
type Config struct {
	KEKName string
	KVMount string
	KVPath  string
}

// DefaultConfig 는 예제에서 쓰는 기본 위치다. policies/*.hcl 의 경로와 맞춰야 한다.
var DefaultConfig = Config{
	KEKName: "member-api-kek",
	KVMount: "secret",
	KVPath:  "member-api/dek",
}

// InitDEK 는 DEK 문서가 없을 때만 v1 을 만든다. 이미 있으면 아무것도 쓰지 않고 false 를 돌려준다.
func InitDEK(ctx context.Context, c *vaultapi.Client, cfg Config) (bool, error) {
	secret, err := c.KVv2(cfg.KVMount).Get(ctx, cfg.KVPath)
	if err == nil {
		if secret.Data == nil {
			return false, fmt.Errorf("%w: %s/%s", ErrDEKDeleted, cfg.KVMount, cfg.KVPath)
		}
		return false, nil
	}
	if !errors.Is(err, vaultapi.ErrSecretNotFound) {
		return false, fmt.Errorf("read DEK %s/%s: %w", cfg.KVMount, cfg.KVPath, err)
	}

	if _, err := RotateDEK(ctx, c, cfg); err != nil {
		return false, err
	}
	return true, nil
}

// RotateDEK 는 Transit 으로 새 DEK 를 발급해 가장 큰 기존 버전의 다음 버전으로 추가하고,
// current_version 을 그 버전으로 바꾼다.
// 기존 버전 필드는 그대로 남겨서 옛 데이터를 계속 복호화할 수 있게 한다.
// 동시에 두 번 실행되면 check-and-set 때문에 뒤에 쓰는 쪽이 실패한다.
func RotateDEK(ctx context.Context, c *vaultapi.Client, cfg Config) (uint8, error) {
	data := map[string]interface{}{}
	wrapped := map[uint8]string{}
	cas := 0

	secret, err := c.KVv2(cfg.KVMount).Get(ctx, cfg.KVPath)
	switch {
	case err == nil:
		if secret.Data == nil {
			return 0, fmt.Errorf("%w: %s/%s", ErrDEKDeleted, cfg.KVMount, cfg.KVPath)
		}
		if _, wrapped, err = parseDEKData(secret.Data); err != nil {
			return 0, err
		}
		for k, v := range secret.Data {
			data[k] = v
		}
		cas = secret.VersionMetadata.Version
	case errors.Is(err, vaultapi.ErrSecretNotFound):
		// 첫 DEK
	default:
		return 0, fmt.Errorf("read DEK %s/%s: %w", cfg.KVMount, cfg.KVPath, err)
	}

	next, err := nextVersion(wrapped)
	if err != nil {
		return 0, err
	}

	resp, err := c.Logical().WriteWithContext(ctx, "transit/datakey/wrapped/"+cfg.KEKName, map[string]interface{}{
		"bits": dekBits,
	})
	if err != nil {
		return 0, fmt.Errorf("generate DEK with %s: %w", cfg.KEKName, err)
	}
	ciphertext, _ := resp.Data["ciphertext"].(string)
	if ciphertext == "" {
		return 0, errors.New("transit datakey returned no ciphertext")
	}

	data[wrappedField(next)] = ciphertext
	data[currentVersionField] = strconv.Itoa(int(next))
	if err := writeDEK(ctx, c, cfg, data, cas); err != nil {
		return 0, err
	}
	return next, nil
}

// LoadKeyring 은 KV 에서 wrap 된 DEK 를 모두 읽고 Transit 으로 unwrap 해서 Keyring 을 만든다.
// 평문 DEK 는 반환된 Keyring 안(메모리)에만 있다.
func LoadKeyring(ctx context.Context, c *vaultapi.Client, cfg Config) (*Keyring, error) {
	secret, err := c.KVv2(cfg.KVMount).Get(ctx, cfg.KVPath)
	if errors.Is(err, vaultapi.ErrSecretNotFound) {
		return nil, fmt.Errorf("%w: %s/%s", ErrDEKNotInitialized, cfg.KVMount, cfg.KVPath)
	}
	if err != nil {
		return nil, fmt.Errorf("read DEK %s/%s: %w", cfg.KVMount, cfg.KVPath, err)
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("%w: %s/%s", ErrDEKDeleted, cfg.KVMount, cfg.KVPath)
	}

	current, wrapped, err := parseDEKData(secret.Data)
	if err != nil {
		return nil, err
	}

	keys := make(map[uint8][]byte, len(wrapped))
	for version, ciphertext := range wrapped {
		key, err := unwrapDEK(ctx, c, cfg.KEKName, ciphertext)
		if err != nil {
			return nil, fmt.Errorf("unwrap DEK v%d: %w", version, err)
		}
		keys[version] = key
	}
	return NewKeyring(current, keys)
}

func unwrapDEK(ctx context.Context, c *vaultapi.Client, kek, ciphertext string) ([]byte, error) {
	resp, err := c.Logical().WriteWithContext(ctx, "transit/decrypt/"+kek, map[string]interface{}{
		"ciphertext": ciphertext,
	})
	if err != nil {
		return nil, err
	}
	encoded, _ := resp.Data["plaintext"].(string)
	if encoded == "" {
		return nil, errors.New("transit decrypt returned no plaintext")
	}
	return base64.StdEncoding.DecodeString(encoded)
}

func writeDEK(ctx context.Context, c *vaultapi.Client, cfg Config, data map[string]interface{}, cas int) error {
	if _, err := c.KVv2(cfg.KVMount).Put(ctx, cfg.KVPath, data, vaultapi.WithCheckAndSet(cas)); err != nil {
		return fmt.Errorf("write DEK %s/%s: %w", cfg.KVMount, cfg.KVPath, err)
	}
	return nil
}

func parseDEKData(data map[string]interface{}) (uint8, map[uint8]string, error) {
	raw, _ := data[currentVersionField].(string)
	current, err := parseVersion(raw)
	if err != nil {
		return 0, nil, fmt.Errorf("%s: %w", currentVersionField, err)
	}

	wrapped := map[uint8]string{}
	for field, value := range data {
		if !strings.HasPrefix(field, wrappedFieldPrefix) {
			continue
		}
		version, err := parseVersion(strings.TrimPrefix(field, wrappedFieldPrefix))
		if err != nil {
			return 0, nil, fmt.Errorf("%s: %w", field, err)
		}
		ciphertext, _ := value.(string)
		if ciphertext == "" {
			return 0, nil, fmt.Errorf("%s is empty", field)
		}
		wrapped[version] = ciphertext
	}

	if _, ok := wrapped[current]; !ok {
		return 0, nil, fmt.Errorf("%s=%d but %s is missing", currentVersionField, current, wrappedField(current))
	}
	return current, wrapped, nil
}

func parseVersion(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 10, 8)
	// "01" 처럼 앞에 0 이 붙은 표기는 wrapped_v1 과 같은 버전으로 겹치므로 거부한다
	if err != nil || v == 0 || strconv.FormatUint(v, 10) != s {
		return 0, fmt.Errorf("invalid DEK version %q (want 1-255)", s)
	}
	return uint8(v), nil
}

// nextVersion 은 이미 있는 DEK 버전 중 가장 큰 값의 다음 버전을 돌려준다.
// current_version 이 되돌려져 있어도 기존 버전 필드를 덮어쓰지 않기 위해서다.
func nextVersion(wrapped map[uint8]string) (uint8, error) {
	var highest uint8
	for version := range wrapped {
		if version > highest {
			highest = version
		}
	}
	if highest == 255 {
		return 0, errors.New("DEK version limit (255) reached")
	}
	return highest + 1, nil
}

func wrappedField(version uint8) string {
	return wrappedFieldPrefix + strconv.Itoa(int(version))
}
