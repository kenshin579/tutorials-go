// member-cli 는 envelope encryption 예제를 실행하는 CLI 다.
//
// 앱 역할(member-api token):        put, get
// migration 역할(member-api-migration token): init-dek, rotate-dek, reencrypt
//
// Vault 주소와 token 은 VAULT_ADDR, VAULT_TOKEN 환경 변수로 받는다.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/kenshin579/tutorials-go/cloud/vault/envelope"
	"github.com/kenshin579/tutorials-go/cloud/vault/member"
)

const usage = `usage: member-cli <command> [flags]

commands:
  init-dek                          create DEK v1 if it does not exist (migration token)
  rotate-dek                        add a new DEK version (migration token)
  reencrypt                         re-encrypt members with the current DEK (migration token)
  put --id ID --email E --phone P   save a member (app token)
  get --id ID                       read a member (app token)

flags:
  --data PATH   member data file (default ./data/members.json)`

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New(usage)
	}
	cmd := args[0]

	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	data := fs.String("data", "./data/members.json", "member data file")
	id := fs.String("id", "", "member id")
	email := fs.String("email", "", "member email")
	phone := fs.String("phone", "", "member phone")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}
	cfg := envelope.DefaultConfig

	switch cmd {
	case "init-dek":
		created, err := envelope.InitDEK(ctx, client, cfg)
		if err != nil {
			return err
		}
		if created {
			fmt.Println("DEK v1 created")
		} else {
			fmt.Println("DEK already exists, skipped")
		}
	case "rotate-dek":
		version, err := envelope.RotateDEK(ctx, client, cfg)
		if err != nil {
			return err
		}
		fmt.Printf("DEK rotated: current_version=%d\n", version)
	case "reencrypt":
		repo, err := openRepository(ctx, client, cfg, *data)
		if err != nil {
			return err
		}
		n, err := repo.Reencrypt()
		if err != nil {
			return err
		}
		fmt.Printf("re-encrypted %d member(s)\n", n)
	case "put":
		if *id == "" || *email == "" || *phone == "" {
			return errors.New("put requires --id, --email and --phone")
		}
		repo, err := openRepository(ctx, client, cfg, *data)
		if err != nil {
			return err
		}
		if err := repo.Put(member.Member{ID: *id, Email: *email, Phone: *phone}); err != nil {
			return err
		}
		fmt.Printf("saved %s\n", *id)
	case "get":
		if *id == "" {
			return errors.New("get requires --id")
		}
		repo, err := openRepository(ctx, client, cfg, *data)
		if err != nil {
			return err
		}
		m, err := repo.Get(*id)
		if err != nil {
			return err
		}
		fmt.Printf("id=%s email=%s phone=%s\n", m.ID, m.Email, m.Phone)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
	}
	return nil
}

func newClient() (*vaultapi.Client, error) {
	cfg := vaultapi.DefaultConfig()
	if cfg.Error != nil {
		return nil, cfg.Error
	}
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	if client.Token() == "" {
		return nil, errors.New("VAULT_TOKEN is not set")
	}
	return client, nil
}

// openRepository 는 Vault 에서 DEK 를 받아 keyring 을 만들고 Repository 를 연다.
// 실제 서비스라면 부팅할 때 한 번 하는 일이다.
func openRepository(ctx context.Context, client *vaultapi.Client, cfg envelope.Config, path string) (*member.Repository, error) {
	k, err := envelope.LoadKeyring(ctx, client, cfg)
	if err != nil {
		return nil, err
	}
	return member.NewRepository(path, k), nil
}
