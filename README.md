# go-auth

JSON Web Token authentication back-end in go

This application aims to provide a simple api to authenticate different
users across different services. Clients pass credentials to this
service and receive a token upon successful authentication. The token is
then used to validate the user for other services in the same realm.
More info on the JWT specification is available in the [ietf
draft](http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html).

**All requests** in which the token is interchanged must be carried out
over an encrypted channel. A malicious third party could easily obtain
the token otherwise, and act on behalf of the victim.

## Build and run the service

Clone and build binary. Expects that the cloned directory is present in
your $GOPATH

    $ git clone https://github.com/aspic/go-auth
    $ cd go-auth/
    $ go get && go build

### Configuration

An example configuration is located in **auth.config.example**. In order
to be able to run go-auth this file must be copied to **auth.config**,
and modified with your credentials.

#### Simple Auth

This scheme is configured as displayed below:

    Auth = simpleAuth // Tells go-auth to use the simpleAuth backend.
    Username = user // Some username
    Password = password // Some password
    Key = key // A key to sign JWTs

Upon authentication go-auth will match username/password from the
request with the configured values. This scheme is most applicable for
testing and initial setup of the application.

### Usage

Run the service, and specify host and port:

    $ ./go-auth -local="localhost:8080"

If you have stock configuration a token can be retrieved by issuing:

    $ curl http://localhost:8080/auth?username=username&password=password

The client has the responsebility to store this token. In subsequent
calls to protected resources the client can present this token to verify
itself.

## Plug into service

An example on how to plug this authtenciation into your go-service is described below. I left out some details for readability. This service will validate the provided token based on its private key (the key corresponding with the key that originially was used to sign the token). 

    // Import client
    import (
        "github.com/aspic/go-auth/client"
        .. other imports
    )
    
    // Setup http handler
    func protectedService(w http.ResponseWriter, r *http.Request) {
    
        // Authenticates based on header, param or cookie
        token := client.AuthByRequest(r, "YOUR APPLICATION KEY")
    
        // A validated token
        if token != nil {
            fmt.Fprintf(w, "Welcome to this protected resource: %s", token.Get("user"))
        } else {
            http.Error(w, "You are not authenticated", http.StatusForbidden)
        }
    }
    
    func main() {
        http.HandleFunc("/protected", protectedService)
        
        .. do stuff
    }
