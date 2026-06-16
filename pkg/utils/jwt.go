package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret") // Ganti dengan ENV di production

type Claims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat token untuk user yang login
func GenerateJWT(userID string, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(80784 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT memverifikasi dan mengurai token
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}

// VerifyToken memverifikasi token
func VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("claims tidak sesuai")
	}

	// log.Printf("\nToken valid: %v | Claims: %+v\n", token.Valid, token.Claims)
	if !ok || !token.Valid || claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		// fmt.Println("tidak valid", token.Valid, claims.ExpiresAt, time.Now())
		// if !ok || !token.Valid || time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	return claims, nil
}
