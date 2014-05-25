/**
 * This package provide types of authentication methods.
 **/
package backend

import (
    "log"
    "database/sql"
    "fmt"
)

var config *Config
var db *sql.DB

type Auth func(username string, password string, realm string) bool

type Config struct {
    Key string      // JWT Secret
    Auth string     // Authentication method
    Username string // Username (for auth or db)
    Password string // Password (for auth or db)
    Host string     // database host
    Database string // database
    Expire int      // Hours until token expire
}

// Reads username and password from the configuration file
func simpleAuth(username string, password string, realm string) bool {
    return username == config.Username && password == config.Password
}

// Auths with postgresql database as back end
func pgAuth(username string, password string, realm string) bool {

    return false
}

func New (conf *Config) Auth {
    var err error
    config = conf

    if config.Auth == "simpleAuth" {
        return simpleAuth
    } else if config.Auth == "pgAuth" {
        // Load database
        props := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", config.Username, config.Password, config.Host, config.Database)
        db, err = sql.Open("postgres", props)
        if err != nil {
            log.Fatal(err)
            return nil
        }
        return pgAuth
    }
    return nil
}
