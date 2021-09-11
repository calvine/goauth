# goauth

An OAuth OIDC server written in Go.

## Current Functionaity

* I am working on it 😜

## TODO List

* Build docker file
* Implement OAuth Autorization Code flow
  * Implement PKCE with Authorization code flow
* Support per app configuration with scopes per app
* Add well known endpoint
* have JWT signing and validation be configuration driven / support (RSA/ ECDSA)

## Notes

* To run mail dev server `docker run --rm -p 1080:1080 -p 1025:1025 -p 8087:8087 maildev/maildev bin/maildev --web 1080 --smtp 1025`
