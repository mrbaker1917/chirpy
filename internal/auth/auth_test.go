package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)
	expiredToken, _ := MakeJWT(userID, "secret", -time.Hour)
	invalidString := "not.a.jwt"

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired token",
			tokenString: expiredToken,
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid string",
			tokenString: invalidString,
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	// Using the Set method (canonicalizes the key)
	headers.Set("Content-Type", "text/json")
	headers.Set("Authorization", "Bearer TOKEN_STRING")
	no_auth_headers := http.Header{}
	no_auth_headers.Set("Content-Type", "text")
	no_bearer := http.Header{}
	no_bearer.Set("Content-Type", "text/html")
	no_bearer.Set("Authorization", "TOKEN_STRING")
	extra_spaces := http.Header{}
	extra_spaces.Set("Content-Type", "text/json")
	extra_spaces.Set("Authorization", "Bearer   MY_TOKEN  ")

	tests := []struct {
		name       string
		hdrs       http.Header
		wantResult string
		wantErr    bool
	}{
		{
			name:       "valid token string",
			hdrs:       headers,
			wantResult: "TOKEN_STRING",
			wantErr:    false,
		},
		{
			name:       "no auth in header",
			hdrs:       no_auth_headers,
			wantResult: "",
			wantErr:    true,
		},
		{
			name:       "no bearer prefix",
			hdrs:       no_bearer,
			wantResult: "",
			wantErr:    true,
		},
		{
			name:       "extra spaces",
			hdrs:       extra_spaces,
			wantResult: "MY_TOKEN",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toke_str, err := GetBearerToken(tt.hdrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if toke_str != tt.wantResult {
				t.Errorf("GetBearerToken() toke_str = %v, want %v", toke_str, tt.wantResult)
			}
		})
	}

}
