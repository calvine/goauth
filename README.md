# goauth

An OAuth OIDC server written in Go.

## Current Functionaity

* I am working on it 😜

## Open Telemetry Use Guidelines

The idea here is to primarily lean on tracing, and add metrics after the fact. Events are added to spans through request processing. each function being traced should add events for the end of the happy path as well as for any error that originates in that function, events shouild not be added to the span for errors who already had events added in other (nested) function calls. if the request fails then the span status should be set to error, but the error details span should not be added again at higher levels of the call stack unless the context makes it relevant.

## TODO List

* Build docker file
* Implement OAuth Autorization Code flow
  * Implement PKCE with Authorization code flow
* Support per app configuration with scopes per app
* Add well known endpoint
* have JWT signing and validation be configuration driven / support (RSA/ ECDSA)

## Notes

* To run mail dev server:
  * `docker run --rm -p 1080:1080 -p 1025:1025 -p 8087:8087 maildev/maildev bin/maildev --web 1080 --smtp 1025`
* Tocreate mongo container for tests or local development:
  * Note the mongo instance to run requires a replica set because we are using transactions.
  * `docker container create --name mongo -p 27017:27017 --env MONGO_INITDB_ROOT_USERNAME=root --env MONGO_INITDB_ROOT_PASSWORD=password --env MONGO_INITDB_DATABASE=testdb mongo:latest --replSet rs0 --bind_ip 0.0.0.0`
