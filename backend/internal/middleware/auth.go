package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserSubKey contextKey = "cognitoSub"

type cognitoJWKS struct {
	Keys []cognitoKey `json:"keys"`
}

type cognitoKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	Use string `json:"use"`
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		sub, err := extractSubFromToken(token)
		if err != nil {
			http.Error(writer, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), UserSubKey, sub)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetUserSub(request *http.Request) (string, bool) {
	sub, ok := request.Context().Value(UserSubKey).(string)
	return sub, ok
}

func extractSubFromToken(tokenString string) (string, error) {
	region := os.Getenv("COGNITO_REGION")
	userPoolId := os.Getenv("COGNITO_USER_POOL_ID")

	jwksUrl := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region,
		userPoolId,
	)

	response, err := http.Get(jwksUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	defer response.Body.Close()

	var jwks cognitoJWKS
	if err := json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return "", fmt.Errorf("failed to decode JWKS: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return parseRSAPublicKey(key.N, key.E)
			}
		}

		return nil, fmt.Errorf("matching key not found")
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub claim not found")
	}

	return sub, nil
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %v", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}
