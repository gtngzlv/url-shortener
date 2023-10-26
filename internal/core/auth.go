package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Claims contains basic jwt.RegisteredClaims and UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const tokenExp = time.Hour * 1
const secretKey = "secret"
const cookieName = "token"

// GetUserToken parses user token from incoming request and returns userID and err if exists
func GetUserToken(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		cookie *http.Cookie
		err    error
	)

	cookie, _ = r.Cookie(cookieName)
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
			return []byte(secretKey), nil
		})
	if err != nil || !token.Valid {
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
		Name:  cookieName,
		Value: token,
		Path:  "/",
	}, nil
}

func generateJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.NewString(),
	})
	return token.SignedString([]byte(secretKey))
}
