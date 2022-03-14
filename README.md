
# Task Management (Backend)
Base On https://github.com/heroku/go-getting-started.git

## Build and Run On Local
```sh
$ go mod tidy
$ go mod vendor
$ go build -o bin/task-management -v .
$ heroku local
```

## Table Schema
```sql
CREATE TABLE public.users (
	id serial4 NOT NULL,
	username varchar NOT NULL,
	"password" varchar NULL,
	nickname varchar NULL,
	CONSTRAINT user_pk PRIMARY KEY (id)
);

CREATE TABLE public.labels (
	id serial4 NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT label_pk PRIMARY KEY (id)
);

CREATE TABLE public.tasks (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	title varchar NOT NULL,
	description varchar NULL,
	start_time timestamptz NULL,
	end_time timestamptz NULL,
	priority int4 NULL,
	state int4 NULL,
	CONSTRAINT task_pk PRIMARY KEY (id),
	CONSTRAINT task_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.task_label_mapping (
	id serial4 NOT NULL,
	task_id int4 NOT NULL,
	label_id int4 NOT NULL,
	CONSTRAINT task_label_mapping_pk PRIMARY KEY (id),
	CONSTRAINT task_label_mapping_fk FOREIGN KEY (task_id) REFERENCES public.tasks(id),
	CONSTRAINT task_label_mapping_fk_1 FOREIGN KEY (label_id) REFERENCES public.labels(id)
);

```

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