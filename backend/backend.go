/**
 * This package provide types of authentication methods.
 **/
package backend

import (
    "log"
    "database/sql"
    _ "github.com/lib/pq"
    "crypto/sha256"
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
    var hash string
    var salt string

    stmt, err := db.Prepare(
        `SELECT u.pw_hash, u.salt FROM identity AS u, realm, inrealm
        WHERE inrealm.id = u.id AND inrealm.realm = realm.id
        AND realm.name = $1 AND u.username = $2`)
    if err != nil {
        log.Print("Failed to execute query: ", err)
        return false
    }
    stmt.QueryRow(realm, username).Scan(&hash, &salt)
    if hash != "" && salt != "" {
        return validPassword(password, salt, hash)
    }
    log.Printf("Unable to authenticate user: %s for realm %s", username, realm)

    return false
}

func validPassword(pw string, salt string, pwHash string) bool {
    pwBytes := []byte(salt + pw)
    hasher := sha256.New()
    hasher.Write(pwBytes)
    sum := fmt.Sprintf("%x", hasher.Sum(nil))
    return sum == pwHash
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
    } else {
        log.Fatal("No backend set, fix configuration")
    }
    return nil
}
