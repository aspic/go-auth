package main

import (
    "flag"
    "fmt"
    "github.com/aspic/go-auth/backend"
    "github.com/aspic/go-auth/client"
    "github.com/dgrijalva/jwt-go"
    "github.com/itkinside/itkconfig"
    "log"
    "net/http"
    "time"
)

var auth backend.Auth
var key string
var config *backend.Config

// Expects username and password, returns token
func authHandler(w http.ResponseWriter, r *http.Request) {

    username := r.FormValue("username")
    password := r.FormValue("password")
    realm := r.FormValue("realm")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Do authentication, store key
    log.Printf("Challenge by '%s' in realm '%s' from IP '%s'", username, realm, r.RemoteAddr)
    success, key := auth(username, password, realm)

    if success {
        token := jwt.New(jwt.GetSigningMethod("HS256"))
        token.Claims["user"] = username
        token.Claims["iss"] = realm
        token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Expire)).Unix()

        tokenString, err := token.SignedString([]byte(key))
        if err == nil {
            log.Print("Success")
            fmt.Fprint(w, tokenString)
        } else {
            log.Print("Error creating token string, ", err)
            http.Error(w, http.StatusText(500), 500)
        }
    } else {
        log.Printf("Invalid username/password/realm combo")
        http.Error(w, http.StatusText(401), 401)
    }
}

// Handler to test token auth
func testHandler(w http.ResponseWriter, r *http.Request) {
    token := client.AuthByRequest(r, config.Key)

    if token != nil {
        fmt.Fprintf(w, "You are: %s", token.Get("user"))
    } else {
        http.Error(w, "You are not authenticated", http.StatusForbidden)
    }
}

func main() {

    // Load config
    config = &backend.Config{Expire: 72}
    configFile := "auth.config"
    err := itkconfig.LoadConfig(configFile, config)
    if err != nil {
        log.Print("Could not read config file ", configFile, err)
    }

    auth = backend.New(config)

    http.HandleFunc("/auth", authHandler)
    http.HandleFunc("/secret", testHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
    flag.Parse()

    if *local != "" {
        err = http.ListenAndServe(*local, nil)
    }
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
