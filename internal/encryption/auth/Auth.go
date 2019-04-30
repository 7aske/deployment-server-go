package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func GenerateToken(secret []byte) string {
	expires := time.Now().Unix() + int64(24*time.Hour)
	type JSTClaims struct {
		Data string `json:"data"`
		jwt.StandardClaims
	}
	// TODO: data
	claims := JSTClaims{
		"",
		jwt.StandardClaims{ExpiresAt: expires, Issuer: "issuer.7aske"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}
func Hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func VerifyCredentials(auser, apass, user, pass string) bool {
	configPassHash := Hash(apass)
	passHash := Hash(pass)
	return configPassHash == passHash && auser == strings.ToLower(user)
}
func VerifyToken(tokenString string, secret []byte) bool {
	if _, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method %v", token.Header["alg"])
		}
		return secret, nil
	}); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}



