package main

import (
    "flag"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/itkinside/itkconfig"
    "github.com/aspic/go-auth/backend"
    "github.com/aspic/go-auth/client"
    "html/template"
    "log"
    "net/http"
    "time"
)

var auther backend.Auther
var config *backend.Config

// Expects username and password, returns token
func authHandler(w http.ResponseWriter, r *http.Request) {

    username := r.FormValue("username")
    password := r.FormValue("password")
    realm := r.FormValue("realm")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Do authentication, store key
    log.Printf("Challenge by '%s' in realm '%s' from IP '%s'", username, realm, r.RemoteAddr)
    success, key := auther.Auth(username, password, realm)

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

// Validate user based on an "arbitrary" token
func idHandler(w http.ResponseWriter, r *http.Request) {
    // Parse username and realm from token, could be forged
    user := client.ParseUser(r)
    if user == nil {
        http.Error(w, "You are not authenticated", http.StatusForbidden)
        return
    }
    // Fetch key based on username and optionally realm
    authInfo := auther.ValidateByUser(user)
    if authInfo.Key == "" {
        http.Error(w, "You are not authenticated", http.StatusForbidden)
        return
    }
    // Validate token based on key. A token is signed by the client application, hence we know if it is forged or not
    token := client.AuthByRequest(r, authInfo.Key)
    if token != nil {
        t, err := template.ParseFiles("templates/id.tpl")
        if err != nil {
            log.Print(err)
        }
        t.Execute(w, authInfo)
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

    auther = backend.New(config)

    http.HandleFunc("/auth", authHandler)
    http.HandleFunc("/id", idHandler)
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
