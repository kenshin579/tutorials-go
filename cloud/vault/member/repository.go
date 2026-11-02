// Package member 는 회원 개인정보(이메일, 전화번호)를 필드 단위로 암호화해 JSON 파일에 저장한다.
package member

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kenshin579/tutorials-go/cloud/vault/envelope"
)

// ErrNotFound 는 ID 에 해당하는 회원이 없을 때 반환된다.
var ErrNotFound = errors.New("member not found")

// Member 는 복호화된 회원 정보다.
type Member struct {
	ID    string
	Email string
	Phone string
}

// record 는 파일에 저장되는 모양이다. Email·Phone 은 base64 로 인코딩한 암호문이다.
type record struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type document struct {
	Members []record `json:"members"`
}

// Repository 는 Keyring 으로 암복호화하며 회원을 JSON 파일에 읽고 쓴다.
type Repository struct {
	path    string
	keyring *envelope.Keyring
}

// NewRepository 는 path 의 파일을 쓰는 Repository 를 만든다. 파일은 처음 Put 할 때 생긴다.
func NewRepository(path string, k *envelope.Keyring) *Repository {
	return &Repository{path: path, keyring: k}
}

// Put 은 Email·Phone 을 암호화해 저장한다. 같은 ID 가 있으면 덮어쓴다.
func (r *Repository) Put(m Member) error {
	email, err := r.seal(m.Email)
	if err != nil {
		return err
	}
	phone, err := r.seal(m.Phone)
	if err != nil {
		return err
	}

	records, err := r.load()
	if err != nil {
		return err
	}
	rec := record{ID: m.ID, Email: email, Phone: phone}
	replaced := false
	for i := range records {
		if records[i].ID == m.ID {
			records[i] = rec
			replaced = true
		}
	}
	if !replaced {
		records = append(records, rec)
	}
	return r.save(records)
}

// Get 은 ID 로 회원을 찾아 복호화한다. 키가 맞지 않으면 에러와 빈 Member 를 돌려준다.
func (r *Repository) Get(id string) (Member, error) {
	records, err := r.load()
	if err != nil {
		return Member{}, err
	}
	for _, rec := range records {
		if rec.ID != id {
			continue
		}
		email, err := r.open(rec.Email)
		if err != nil {
			return Member{}, fmt.Errorf("decrypt email of %s: %w", id, err)
		}
		phone, err := r.open(rec.Phone)
		if err != nil {
			return Member{}, fmt.Errorf("decrypt phone of %s: %w", id, err)
		}
		return Member{ID: id, Email: email, Phone: phone}, nil
	}
	return Member{}, fmt.Errorf("%w: %s", ErrNotFound, id)
}

// Reencrypt 는 current 가 아닌 DEK 버전으로 암호화된 필드를 current 버전으로 다시 암호화하고,
// 바뀐 회원 수를 돌려준다. DEK 를 회전한 뒤 옛 버전을 정리할 때 쓴다.
func (r *Repository) Reencrypt() (int, error) {
	records, err := r.load()
	if err != nil {
		return 0, err
	}

	changed := 0
	for i := range records {
		emailChanged, err := r.reseal(&records[i].Email)
		if err != nil {
			return 0, fmt.Errorf("re-encrypt email of %s: %w", records[i].ID, err)
		}
		phoneChanged, err := r.reseal(&records[i].Phone)
		if err != nil {
			return 0, fmt.Errorf("re-encrypt phone of %s: %w", records[i].ID, err)
		}
		if emailChanged || phoneChanged {
			changed++
		}
	}

	if changed == 0 {
		return 0, nil
	}
	return changed, r.save(records)
}

func (r *Repository) seal(plaintext string) (string, error) {
	blob, err := r.keyring.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(blob), nil
}

func (r *Repository) open(encoded string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	plaintext, err := r.keyring.Decrypt(blob)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// reseal 은 필드가 current 버전이 아니면 다시 암호화하고 바꿨는지를 돌려준다.
func (r *Repository) reseal(field *string) (bool, error) {
	blob, err := base64.StdEncoding.DecodeString(*field)
	if err != nil {
		return false, err
	}
	version, err := envelope.VersionOf(blob)
	if err != nil {
		return false, err
	}
	if version == r.keyring.Current() {
		return false, nil
	}

	plaintext, err := r.open(*field)
	if err != nil {
		return false, err
	}
	sealed, err := r.seal(plaintext)
	if err != nil {
		return false, err
	}
	*field = sealed
	return true, nil
}

func (r *Repository) load() ([]record, error) {
	raw, err := os.ReadFile(r.path)
	if errors.Is(err, fs.ErrNotExist) || (err == nil && len(raw) == 0) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var doc document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", r.path, err)
	}
	return doc.Members, nil
}

// save 는 임시 파일에 쓴 뒤 rename 해서, 쓰다가 실패해도 기존 파일이 깨지지 않게 한다.
func (r *Repository) save(records []record) error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(document{Members: records}, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".members-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(raw, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), r.path)
}
