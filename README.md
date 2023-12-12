
# Task Management (Backend)
Base On https://github.com/heroku/go-getting-started.git

## Build and Run On Local
```sh
$ go mod tidy
$ go mod vendor
$ go build -o bin/task-management -v .
$ bin/task-management
```

## Environment variables
|Name|Description|
|---|---|
|DATABASE_URL|Database URL|
|PORT|Port Number|
|JWTSIGNKEY|JWT sign key string|
|PASSWORDSALT|password salt string|
|GOOGLE_OAUTH2_CLIENTID|Google OAuth2 ClientID|
|GOOGLE_OAUTH2_CLIENTSECRET|Google OAuth2 ClientSecret|

## APIs
||||
|--|--|--|
|/auth|/login||
||/googleOAuth||
|/apis|/labels|/list|
|||/create|
|||/update|
|||/delete|
||/tasks|/list|
|||/create|
|||/update|
|||/delete|
||/user|/getNickName|
