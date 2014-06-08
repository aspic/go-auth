package common

// Key used for authentication, realms user belongs to
type AuthInfo struct {
    User *User
    Key string
    Realms []*Realm
}

type User struct {
    Realm string `json:"iss"`
    Username string `json:"user"`
    Expire int `json:"exp"`
}

type Realm struct {
    Name string
}
