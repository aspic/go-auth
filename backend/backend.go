/**
 * This package provide types of authentication methods.
 **/
package backend

import (
    "log"
    "database/sql"
    _ "github.com/lib/pq"
    "crypto/sha256"
    "github.com/aspic/go-auth/common"
    "fmt"
)

var config *Config
var db *sql.DB

type Auther interface {
    Auth(username string, password string, realm string) (bool, string)
    ValidateByUser(user *common.User) *common.AuthInfo
}

type Simple struct {}

type Postgres struct {}

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
func (s *Simple) Auth(username string, password string, realm string) (bool, string) {
    return username == config.Username && password == config.Password, config.Key
}

// Auths with postgresql database as back end
func (p *Postgres) Auth(username string, password string, realm string) (bool, string) {
    var hash string
    var salt string
    var key string

    stmt, err := db.Prepare(
        `SELECT u.pw_hash, u.salt, realm.key FROM identity AS u, realm, inrealm
        WHERE inrealm.id = u.id AND inrealm.realm = realm.id
        AND realm.name = $1 AND u.username = $2`)
    if err != nil {
        log.Print("Failed to execute query: ", err)
        return false, ""
    }
    stmt.QueryRow(realm, username).Scan(&hash, &salt, &key)
    log.Print("Got key: ", key)
    if hash != "" && salt != "" {
        return validPassword(password, salt, hash), key
    }
    log.Printf("Unable to authenticate user: %s for realm %s", username, realm)

    return false, ""
}

func (s *Simple) ValidateByUser(user *common.User) *common.AuthInfo {
    if user.Username == config.Username {
        return &common.AuthInfo{Key: config.Key, User: user}
    }
    return nil
}

func (s *Postgres) ValidateByUser(user *common.User) *common.AuthInfo {
    stmt, err := db.Prepare(
        `SELECT r.name FROM identity AS i, realm AS r, inrealm
         WHERE inrealm.id = i.id AND r.id = inrealm.realm AND i.username = $1`)
    if err != nil {
        log.Fatal(err)
    }
    rows, err := stmt.Query(user.Username)
    realms := make([]*common.Realm, 0)
    var key string
    for rows.Next() {
        var realmName string
        var realmKey string
        if err := rows.Scan(&realmName); err != nil {
            log.Fatal(err)
        }
        if realmName == user.Realm {
            key = realmKey
        }
        realms = append(realms, &common.Realm{Name: realmName})
    }
    // Got valid realm and matching key
    if len(realms) > 0 && key != "" {
        return &common.AuthInfo{Key: key, User: user, Realms: realms}
    }

    return nil
}

func validPassword(pw string, salt string, pwHash string) bool {
    pwBytes := []byte(salt + pw)
    hasher := sha256.New()
    hasher.Write(pwBytes)
    sum := fmt.Sprintf("%x", hasher.Sum(nil))
    return sum == pwHash
}

func New (conf *Config) Auther {
    var err error
    config = conf

    if config.Auth == "simpleAuth" {
        return &Simple{}
    } else if config.Auth == "pgAuth" {
        // Load database
        props := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", config.Username, config.Password, config.Host, config.Database)
        db, err = sql.Open("postgres", props)
        if err != nil {
            log.Fatal(err)
            return nil
        }
        return &Postgres{}
    } else {
        log.Fatal("No backend set, fix configuration")
    }
    return nil
}
