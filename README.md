
# Task Management (Backend)
Base On https://github.com/heroku/go-getting-started.git

## Build and Run On Local
```sh
$ go mod tidy
$ go mod vendor
$ go build -o bin/task-management -v .
$ heroku local
```

## Environment variables
|Name|Description|
|---|---|
|DATABASE_URL|Database URL|
|JWTSIGNKEY|JWT sign key string|
|PASSWORDSALT|password salt string|

## APIs
|||
|--|--|
|/labels|/list|
||/create|
||/update|
||/delete|
|/tasks|/list|
||/create|
||/update|
||/delete|
|/auth|/login|
|/user|/getNickName|
