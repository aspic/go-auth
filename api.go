package main

import (
    "log"
    "fmt"
    "net/http"
    "flag"
    "github.com/dgrijalva/jwt-go"
    "github.com/itkinside/itkconfig"
    "github.com/aspic/go-auth/client"
)

var auth Auth
var config *Config

type Config struct {
    Key string      // JWT Secret
    Auth string     // Authentication method
    Username string // Username (for auth or db)
    Password string // Password (for auth or db)
    Host string     // database host
    Database string // database
}

type Auth func(username string, password string, realm string) bool

// Reads username and password from the configuration file
func simpleAuth(username string, password string, realm string) bool {
    return username == config.Username && password == config.Password
}

// Expects username and password, returns token
func authHandler(w http.ResponseWriter, r *http.Request) {

    username := r.FormValue("username")
    password := r.FormValue("password")
    realm := r.FormValue("realm")

    // Do authentication
    if auth(username, password, realm) {
        token := jwt.New(jwt.GetSigningMethod("HS256"))
        token.Claims["user"] = username
        token.Claims["iss"] = realm

        tokenString, err := token.SignedString([]byte(config.Key))
        if err == nil {
            fmt.Fprint(w, tokenString)
        } else {
            log.Print("Error creating token string, ", err)
            http.Error(w, http.StatusText(500), 500)
        }
    } else {
        http.Error(w, http.StatusText(401), 401)
    }
}

// Handler to test token auth
func testHandler(w http.ResponseWriter, r *http.Request) {
    token := client.AuthWithCookie(r, config.Key)

    if token != nil {
        fmt.Fprintf(w, "You are: %s", token.Get("user"))
    } else {
        http.Error(w, "You are not authenticated", http.StatusForbidden)
    }
}

func main() {

    // Load config
    config = &Config{}
    configFile := "auth.config"
    err := itkconfig.LoadConfig(configFile, config)
    if err != nil {
        log.Print("Could not read config file ", configFile, err)
    }

    if config.Auth == "simpleAuth" {
        auth = simpleAuth
    }

    http.HandleFunc("/auth", authHandler)
    http.HandleFunc("/secret", testHandler)

    var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
    flag.Parse()

    if *local != "" {
        err = http.ListenAndServe(*local, nil)
    }
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
