package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Inisialisasi seed hanya sekali
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateInvitationCode membuat kode undangan seperti RT05-ABC123
func GenerateInvitationCode(rt string) string {
	randomPart := generateRandomString(6)
	return fmt.Sprintf("%s-%s", strings.ToUpper(rt), randomPart)
}

// generateRandomString menghasilkan string acak huruf kapital & angka
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
