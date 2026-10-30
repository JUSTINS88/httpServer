package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerTokenCases(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{
			name:   "header valid",
			header: "Bearer test-token",
			want:   "test-token",
		},
		{
			name:    "header tidak ada",
			wantErr: true,
		},
		{
			name:    "skema salah",
			header:  "Basic test-token",
			wantErr: true,
		},
		{
			name:    "token tidak ada",
			header:  "Bearer",
			wantErr: true,
		},
		{
			name:    "token tidak ada",
			header:  "Bearer ",
			wantErr: true,
		},
		{
			name:    "token tidak ada",
			header:  "BearerFake",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetBearerToken(headers)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("token = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJWTLifecycle(t *testing.T) {
	const secret = "secret-khusus-testing"
	userID := uuid.New()

	tests := []struct {
		name           string
		expiresIn      time.Duration
		validateSecret string
		wantErr        bool
	}{
		{
			name:           "token valid",
			expiresIn:      time.Hour,
			validateSecret: secret,
		},
		{
			name:           "token kedaluwarsa",
			expiresIn:      -time.Hour,
			validateSecret: secret,
			wantErr:        true,
		},
		{
			name:           "secret salah",
			expiresIn:      time.Hour,
			validateSecret: "secret-berbeda",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(userID, secret, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT gagal: %v", err)
			}

			gotID, err := ValidateJWT(token, tt.validateSecret)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && gotID != userID {
				t.Errorf("user ID = %v, want %v", gotID, userID)
			}
		})
	}
}
