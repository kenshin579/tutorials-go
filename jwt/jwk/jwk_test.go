package jwk

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/MicahParks/jwkset"
	"github.com/stretchr/testify/suite"
)

type jwksTestSuite struct {
	suite.Suite
	publicKey string
	ctx       context.Context
	jwkCache  jwkset.Storage
}

func TestJwtUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(jwksTestSuite))
}

func (suite *jwksTestSuite) SetupSuite() {
	publicKeyPath := "../pem/public.key"
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	suite.NoError(err)

	suite.ctx = context.TODO()
	suite.jwkCache = jwkset.NewMemoryStorage()
	suite.publicKey = string(publicKeyBytes)
}

func (suite *jwksTestSuite) Test_GenerateJWKS() {
	suite.Run("JWKS의 값으로 Public Key와 같은지 확인한다", func() {
		jwks, err := suite.generateJWKS()
		suite.NoError(err)

		// Write jwks as json file
		file, err := os.Create("jwks.json")
		suite.NoError(err)
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(jwks)
		suite.NoError(err)

		suite.Equal(jwkset.UseSig, jwks.Keys[0].USE)
		jwk, err := jwkset.NewJWKFromMarshal(jwks.Keys[0], jwkset.JWKMarshalOptions{}, jwkset.JWKValidateOptions{})

		originalPubKey, err := parsePublicKey(suite.publicKey)
		suite.NoError(err)
		suite.Equal(originalPubKey, jwk.Key())
	})
}

func (suite *jwksTestSuite) generateJWKS() (jwkset.JWKSMarshal, error) {
	pubKey, err := parsePublicKey(suite.publicKey)
	if err != nil {
		return jwkset.JWKSMarshal{}, err
	}

	jwk, err := jwkset.NewJWKFromKey(pubKey, jwkset.JWKOptions{
		Marshal: jwkset.JWKMarshalOptions{
			Private: false,
		},
		Metadata: jwkset.JWKMetadataOptions{
			USE: jwkset.UseSig,
		},
	})
	if err != nil {
		return jwkset.JWKSMarshal{}, err
	}

	if err := suite.jwkCache.KeyWrite(suite.ctx, jwk); err != nil {
		return jwkset.JWKSMarshal{}, err
	}

	jwkMarshal, err := suite.jwkCache.Marshal(suite.ctx)
	if err != nil {
		return jwkset.JWKSMarshal{}, err
	}
	return jwkMarshal, nil
}

func parsePublicKey(publicKeyPEM string) (any, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	return pubKey, nil
}
