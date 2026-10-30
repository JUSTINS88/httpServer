package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatal(err)
	}

	if gotID != userID {
		t.Errorf("got %v, want %v", gotID, userID)
	}
}
