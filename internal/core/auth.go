package core

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TokenExp = time.Hour * 1
const SecretKey = "secret"
const CookieName = "token"

func GetUserToken(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		cookie *http.Cookie
		err    error
	)

	cookie, _ = r.Cookie(CookieName)
	if cookie == nil {
		cookie, err = generateCookie()
		if err != nil {
			return "", fmt.Errorf("GetUserToken: failed to generate cookie, %s", err)
		}
		http.SetCookie(w, cookie)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, err
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		cookie, err = generateCookie()
		http.SetCookie(w, cookie)
	}

	if !token.Valid {
		cookie, err = generateCookie()
		http.SetCookie(w, cookie)
	}
	return claims.UserID, nil

}

func generateCookie() (*http.Cookie, error) {
	token, err := generateJWTString()
	if err != nil {
		return nil, fmt.Errorf("generateCookie: failed to generate, %s", err)
	}
	return &http.Cookie{
		Name:  CookieName,
		Value: token,
		Path:  "/",
	}, nil
}

func generateJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: uuid.NewString(),
	})
	return token.SignedString([]byte(SecretKey))
}
