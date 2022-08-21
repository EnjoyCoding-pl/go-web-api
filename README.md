# Description #

This is an example of clean architecture in GO web api. The main goals of this project is learning how to:
- handle http requests
- structure code 
- connect application to postgres database 
- add open telemetry with jaeger provider
- authorize users

# Run tests #

`go test ./...`

# Run web api with default configuration #

`go run ./cmd`

# Configuration options #
- address - Api address default 127.0.0.1:4400
- jaeger - Address to jaeger instance
- jwt-issuer - JWT iss claim
- jwt-secret - Secret for signing jwt
- timeout - Api timeout
- connection-string - Connection string to postgres database

# Run web api with custom configuration #

`go run ./cmd --address "127.0.0.1:4400" --jaeger "http://localhost:14268/api/traces" --jwt-issuer "http://127.0.0.1:4400" --jwt-secret "secret_placeholder" --timeout 30 --connection-string "host=127.0.0.1 user=postgres password=postgres dbname=go-web-app port=5432 sslmode=disable" `
