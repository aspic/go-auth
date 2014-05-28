package client

import (
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "log"
)

const (
    TOKEN_FORM = "token"
    TOKEN_HEADER = "x-access-token"
    TOKEN_COOKIE = "token"
)

type Token struct {
    token *jwt.Token
}

/**
 * Looks for query param 'token', header 'x-access-token' and cookie 'token'
 * in that specific order
 */
func AuthByRequest(r *http.Request, key string) *Token {
    token := AuthWithKey(r.FormValue(TOKEN_FORM), key)
    if token != nil {
        return token
    }
    token = AuthWithKey(r.Header.Get(TOKEN_HEADER), key)
    if token != nil {
        return token
    }
    cookie, err := r.Cookie(TOKEN_COOKIE)
    if err == nil && cookie != nil {
        return AuthWithKey(cookie.Value, key)
    }
    log.Print("Unable to read cookie, ", err)
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
