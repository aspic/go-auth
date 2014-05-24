package main

import (
    "log"
    "fmt"
    "net/http"
    "flag"
    "github.com/dgrijalva/jwt-go"
    "github.com/aspic/go-auth/client"
)

const (
    KEY = "LOL"
)

// Expects username and password, returns token
func authHandler(w http.ResponseWriter, r *http.Request) {

    username := r.FormValue("username")
    password := r.FormValue("password")
    authString := auth(username, password)

    if authString != "" {
        fmt.Fprint(w, authString)
    } else {
        http.Error(w, http.StatusText(403), 403)
    }
}

// Handler to test token auth
func testHandler(w http.ResponseWriter, r *http.Request) {
    token := client.AuthWithCookie(r, KEY)

    if token != nil {
        fmt.Fprintf(w, "You are: %s", token.Get("user"))
    } else {
        http.Redirect(w, r, "/js/login.html", http.StatusMovedPermanently)
    }
}

// Authenticate user by some back end
func auth(username string, password string) (token string) {
    if username == "test" && password == "test" {
        token := jwt.New(jwt.GetSigningMethod("HS256"))
        token.Claims["user"] = username
        token.Claims["iss"] = "mehl"

        tokenString, err := token.SignedString([]byte(KEY))
        if err == nil {
            return tokenString
        } else {
            log.Print("Error creating token string, ", err)
        }
    }
    log.Printf("Could not authenticate user: %s", username)
    return ""
}


func main() {

    http.HandleFunc("/secret", testHandler)
    http.HandleFunc("/auth", authHandler)
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))

    var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
    flag.Parse()

    var err error
    if *local != "" {
        err = http.ListenAndServe(*local, nil)
    }
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
