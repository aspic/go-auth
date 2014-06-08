package client

import (
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "github.com/aspic/go-auth/common"
    "log"
    "net/http"
    "strings"
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
 * in that order
 */
func AuthByRequest(r *http.Request, key string) *Token {
    tokenString := getJwtString(r)
    if tokenString != "" {
        return AuthWithKey(tokenString, key)
    }
    return nil
}


func getJwtString(r *http.Request) string {
    tokenString := r.FormValue(TOKEN_FORM)
    if tokenString != "" {
        return tokenString
    }
    tokenString = r.Header.Get(TOKEN_HEADER)
    if tokenString != "" {
        return tokenString
    }
    cookie, err := r.Cookie(TOKEN_COOKIE)
    if err == nil && cookie != nil {
        return cookie.Value
    }
    return ""
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

// Parses part of the token to a User struct
func ParseUser(r *http.Request) *common.User {
    tokenString := getJwtString(r)
    slices := strings.Split(tokenString, ".")
    if len(slices) != 3 {
        log.Print("Can't decode tokenstring: ", tokenString)
        return nil
    }
    bytes, err := jwt.DecodeSegment([1])
    if err != nil {
        log.Print("Unable to decode segment: ", err)
        return nil
    }
    user := &common.User{}
    err = json.Unmarshal(bytes, user)
    if err != nil {
        log.Print("Unable to decode json: ", err)
        return nil
    }
    return user
}
