package client

import (
    "testing"
    "net/http"
    "log"
)

func TestRequestWithHeaderToken_shouldReturnToken(t *testing.T) {
    stringToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MDE1NzEwNDEsImlzcyI6IiIsInVzZXIiOiJyb290In0.yGhhyKdM4trCcTXFmHn--O3Z7XrFxyrK1jVdvEC6jro"
    req, err := http.NewRequest("GET", "http://example.com/foo", nil)
    req.Header.Set("x-access-token", stringToken)
    if err != nil {
        log.Fatal(err)
    }
    token := AuthByRequest(req, "thissecretkey")
    if token == nil {
        t.Errorf("Token not set")
    }
}
