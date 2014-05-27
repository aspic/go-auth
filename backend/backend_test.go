package backend

import "testing"

func TestVerifyPassword(t *testing.T) {
    salt := "somesalt"
    pw := "hemmeleg"
    hash := "6561c49480e895229db8d5712b88e90c71e77116fa3e47b411ebb3d01fa58132"

    if !verifyPassword(pw, salt, hash) {
        t.Errorf("Password could not be verified")
    }
}
