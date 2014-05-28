package client

import (
    "testing"
    "net/http"
    "log"
    "fmt"
    "time"
)

// Pre-generated token, username: root, password: secret, signing-key: thissecretkey
var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MDE1NzEwNDEsImlzcyI6IiIsInVzZXIiOiJyb290In0.yGhhyKdM4trCcTXFmHn--O3Z7XrFxyrK1jVdvEC6jro"

func TestRequestWithHeaderToken_shouldReturnToken(t *testing.T) {
    req, err := http.NewRequest("GET", "http://auth-server.com/foo", nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("x-access-token", tokenString)

    token := AuthByRequest(req, "thissecretkey")
    if token == nil {
        t.Errorf("Token not set")
    }
}

func TestRequestWithRequestToken_shouldReturnToken(t *testing.T) {
    req, err := http.NewRequest("GET", fmt.Sprintf("http://auth-server/?token=%s", tokenString), nil)
    if err != nil {
        log.Fatal(err)
    }

    token := AuthByRequest(req, "thissecretkey")
    if token == nil {
        t.Errorf("Token not set")
    }
}

func TestRequestWithCookieToken_shouldReturnToken(t *testing.T) {
    req, err := http.NewRequest("GET", "http://auth-server/", nil)
    if err != nil {
        log.Fatal(err)
    }
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{
        "token", tokenString, "/", "http://auth-server/", expire, expire.Format(time.UnixDate), 86400,
        true, true, "", []string{}}
    req.AddCookie(&cookie)

    token := AuthByRequest(req, "thissecretkey")
    if token == nil {
        t.Errorf("Token not set")
    }
}
