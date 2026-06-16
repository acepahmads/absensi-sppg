// password.go placeholder
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword meng-hash password plaintext
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword membandingkan hash dengan password plaintext
func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// fmt.Println("Compare error:", err)
	return err == nil
}
