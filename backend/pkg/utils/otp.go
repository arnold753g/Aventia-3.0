package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// GenerateOTP6 genera un código OTP de 6 dígitos aleatorio.
func GenerateOTP6() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// HashOTP genera un hash HMAC-SHA256 del código OTP.
// purpose: "email_verify" o "password_reset".
func HashOTP(purpose, email, code string) string {
	secret := os.Getenv("OTP_SECRET")
	message := fmt.Sprintf("%s:%s:%s", purpose, strings.ToLower(email), code)

	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}

// VerifyOTP verifica que el código coincide con el hash almacenado.
func VerifyOTP(purpose, email, code, storedHash string) bool {
	computedHash := HashOTP(purpose, email, code)
	return hmac.Equal([]byte(computedHash), []byte(storedHash))
}
