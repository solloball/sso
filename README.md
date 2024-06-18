# SSO service
SSO service with grpc on go
There are three services
1. Auth service (in testing)
  There are three handlers:
  1. IsAdmin
  2. Login
  3. Register
5. Permision service (coming soon)
6. User info service (coming soon)

It uses sqlite for database. There is a migrator, for use it exec: "make"

# How to build and run in production mode
To run:
```sh
go run cmd/sso/main.go
```
docker coming soon
