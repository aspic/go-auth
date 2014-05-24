# go-auth

JSON Web Token authentication back-end in go

This application aims to provide a simple api for authentication different users across different services. Clients pass credentials to this service and receives a token upon success. This token is used for further authentication agains other services in the same realm. More info on the JWT specification is available in the [ietf draft](http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html).

Notice that **all requests** in which the token is interchanged must be carried out over an encrypted channel. A malicious third party could easily obtain the token otherwise, and act on behalf of the victim.
