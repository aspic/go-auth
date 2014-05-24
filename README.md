# go-auth

JSON Web Token authentication back-end in go

This application aims to provide a simple api to authenticate different users across different services. Clients pass credentials to this service and receive a token upon successful authentication. The token is then used to validate the user for other services in the same realm. More info on the JWT specification is available in the [ietf draft](http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html).

**All requests** in which the token is interchanged must be carried out over an encrypted channel. A malicious third party could easily obtain the token otherwise, and act on behalf of the victim.
