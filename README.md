
# go-getting-started

A barebones Go app, which can easily be deployed to Heroku.

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ git clone https://github.com/heroku/go-getting-started.git
$ cd go-getting-started
$ go build -o bin/go-getting-started -v . # or `go build -o bin/go-getting-started.exe -v .` in git bash
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku main
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)

## TaskManagement
```sh
$ go mod tidy
$ go mod vendor
$ go build -o bin/task-management -v .
$ heroku local
```

## Table Schema
task

id
user_id
title
description
start_time
end_time
priority
state

label
id
name

task_label_mapping
id
task_id
label_id

user
id
username
password
nickname

```sql
CREATE TABLE public."user" (
	id serial NOT NULL,
	username varchar NOT NULL,
	"password" varchar NULL,
	nickname varchar NULL,
	CONSTRAINT user_pk PRIMARY KEY (id)
);

CREATE TABLE public."label" (
	id serial4 NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT label_pk PRIMARY KEY (id)
);

CREATE TABLE public.task (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	title varchar NOT NULL,
	description varchar NULL,
	start_time timestamptz NULL,
	end_time timestamptz NULL,
	priority int4 NULL,
	state int4 NULL,
	CONSTRAINT task_pk PRIMARY KEY (id),
	CONSTRAINT task_fk FOREIGN KEY (user_id) REFERENCES public."user"(id)
);

CREATE TABLE public.task_label_mapping (
	id serial4 NOT NULL,
	task_id int4 NOT NULL,
	label_id int4 NOT NULL,
	CONSTRAINT task_label_mapping_pk PRIMARY KEY (id),
	CONSTRAINT task_label_mapping_fk FOREIGN KEY (task_id) REFERENCES public.task(id),
	CONSTRAINT task_label_mapping_fk_1 FOREIGN KEY (label_id) REFERENCES public."label"(id)
);

```