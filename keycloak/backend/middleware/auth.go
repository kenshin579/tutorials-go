package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWKSResponse represents the JWKS response from Keycloak
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// KeycloakClaims represents the claims in Keycloak JWT token
type KeycloakClaims struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	jwt.RegisteredClaims
}

// JWTMiddleware creates a middleware for JWT token validation
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse token without verification first to get the kid
			token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Get kid from token header
				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("kid not found in token header")
				}

				// Get public key from Keycloak JWKS endpoint
				publicKey, err := getPublicKeyFromKeycloak(kid)
				if err != nil {
					return nil, fmt.Errorf("failed to get public key: %v", err)
				}

				return publicKey, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Extract claims and store in context
			if claims, ok := token.Claims.(*KeycloakClaims); ok {
				c.Set("user", claims)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			return next(c)
		}
	}
}

// getPublicKeyFromKeycloak retrieves the public key from Keycloak JWKS endpoint
func getPublicKeyFromKeycloak(kid string) (*rsa.PublicKey, error) {
	// Keycloak JWKS URL
	jwksURL := "http://localhost:8080/realms/myrealm/protocol/openid-connect/certs"

	// Fetch JWKS
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %v", err)
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Kty == "RSA" {
			return jwkToRSAPublicKey(key)
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

// jwkToRSAPublicKey converts JWK to RSA public key
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode n (modulus)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	// Decode e (exponent)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Convert to big.Int
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}
