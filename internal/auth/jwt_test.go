package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestGenerateRefreshToken(t *testing.T) {
	token1, err := GenerateRefreshToken()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Refresh Token: %s", token1)
}

func TestHashRefreshToken(t *testing.T) {
	hasher := sha256.New()
	hasher.Write([]byte("8UMBun-uTzNTmFoPErpxvqIZeI3UuFazIA3bjwp0S7w="))

	t.Logf("Hashed Refresh Token: %s", hex.EncodeToString(hasher.Sum(nil)))
}
