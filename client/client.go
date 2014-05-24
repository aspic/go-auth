package client

import (
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "log"
)

type Token struct {
    token *jwt.Token
}

// Validates token with a request cookie
func AuthWithCookie(r *http.Request, key string) *Token {
    cookie, err := r.Cookie("token")
    if err == nil && cookie != nil {
        return AuthWithKey(cookie.Value, key)
    }
    log.Print("Unable to read cookie")
    return nil
}

// Validates token with the provided key.
func AuthWithKey(token string, key string) *Token {
    return Auth(token, func(token *jwt.Token) ([]byte, error) {
        return []byte(key), nil
    })
}

// Validates token with the provided KeyFunc.
func Auth(tokenString string, keyFunc jwt.Keyfunc) *Token {
    token, err := jwt.Parse(tokenString, keyFunc)

    if err == nil && token.Valid {
        return &Token{token: token}
    }
    return nil
}

// Gets the value of the claim key
func (token *Token) Get(key string) string {
    return token.token.Claims[key].(string)
}
