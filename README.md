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

The service is started by executing the binary with hostname and port
specified as an argument.

    $ ./go-auth -local="localhost:8080"

## Usage

A token is retrieved by authenticating with the /auth endpoint:

    $ curl http://localhost:8080/auth?username=foo&password=bar

The client has to store this token, and present it when requesting
services that are protected by the go-auth scheme.
