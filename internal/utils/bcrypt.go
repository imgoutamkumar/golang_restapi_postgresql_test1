package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func HashOtp(otp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckOtpHash(hashedOtp, plainOtp string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedOtp), []byte(plainOtp))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, err
	}
	return true, nil
}

// GenerateSecureToken creates a URL-safe random string
func GenerateSecureToken() (string, error) {
	// 32 bytes of entropy yields a ~43 character base64 string
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Use URLEncoding so it is safe to put directly into a web URL
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
