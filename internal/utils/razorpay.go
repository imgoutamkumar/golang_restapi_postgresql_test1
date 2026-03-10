package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func VerifySignature(body string, signature string) bool {

	secret := os.Getenv("RAZORPAY_SECRET")

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))

	expected := hex.EncodeToString(h.Sum(nil))

	return expected == signature
}
