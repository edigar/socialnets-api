package authentication

import (
	"github.com/edigar/socialnets-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	userId := "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"
	token, err := CreateToken(userId)
	if err != nil {
		t.Errorf("CreateToken should not return an error for a valid uint64: %v", err)
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return config.SecretKey, nil
	})

	if err != nil {
		t.Errorf("ParseWithClaims should not return an error: %v", err)
	}
	if !parsedToken.Valid {
		t.Errorf("Parsed token should be valid")
	}
	if claims["authorized"] != true {
		t.Errorf("Token must have an authorized claim set to true")
	}
	if claims["exp"] == nil {
		t.Errorf("Token should have an exp claim not nil")
	}

	tokenUserId := claims["userId"]
	if err != nil {
		t.Errorf("ParseUint of userId claims should not return an error: %v", err)
	}
	if userId != tokenUserId {
		t.Errorf("Token should have a correct userId")
	}
}

func TestTokenValidateValidToken(t *testing.T) {
	userId := "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"
	token, err := CreateToken(userId)
	if err != nil {
		t.Errorf("CreateToken should not return an error for a valid uint64: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	err = TokenValidate(request)
	if err != nil {
		t.Errorf("TokenValidate should not return an error: %v", err)
	}
}

func TestTokenValidateInvalidToken(t *testing.T) {
	token := "invalid-token"

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	err := TokenValidate(request)

	if err == nil {
		t.Errorf("TokenValidate should return an error for an invalid token")
	}
}

func TestExtractUserIdValidToken(t *testing.T) {
	userId := "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"
	token, err := CreateToken(userId)
	if err != nil {
		t.Errorf("CreateToken should not return an error")
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	extractedUserId, err := ExtractUserId(request)

	if err != nil {
		t.Errorf("ExtractUserId should not return an error")
	}
	if userId != extractedUserId {
		t.Errorf(
			"Extracted user ID should match the token's user ID. UserId: %v. extractedUserId: %v",
			userId,
			extractedUserId,
		)
	}
}

func TestExtractUserIdInvalidToken(t *testing.T) {
	token := "invalid-token"

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	_, err := ExtractUserId(request)
	if err == nil {
		t.Errorf("ExtractUserId should return an error for an invalid token")
	}
}

func TestTokenValidateRejectsNonHMACSigningMethod(t *testing.T) {
	userId := "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Hour).Unix(),
		"userId":     userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("Signing a none-method token should not fail: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)

	if err = TokenValidate(request); err == nil {
		t.Errorf("TokenValidate should reject a token signed with a non-HMAC method (algorithm confusion)")
	}
	if _, err = ExtractUserId(request); err == nil {
		t.Errorf("ExtractUserId should reject a token signed with a non-HMAC method (algorithm confusion)")
	}
}

func TestTokenValidateRejectsExpiredToken(t *testing.T) {
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(-time.Hour).Unix(),
		"userId":     "eedf21bf-dde8-4c85-b50b-89a1cba87c2e",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		t.Fatalf("Signing the token should not fail: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)

	if err = TokenValidate(request); err == nil {
		t.Errorf("TokenValidate should reject an expired token")
	}
}

func TestTokenValidateRejectsMissingOrMalformedHeader(t *testing.T) {
	t.Run("missing Authorization header", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		if err := TokenValidate(request); err == nil {
			t.Errorf("TokenValidate should return an error when the Authorization header is missing")
		}
	})

	t.Run("Authorization header without bearer token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Authorization", "Bearer")

		if err := TokenValidate(request); err == nil {
			t.Errorf("TokenValidate should return an error when the Authorization header has no bearer token")
		}
	})
}
