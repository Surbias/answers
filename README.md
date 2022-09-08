# Answers API

## General usage:
- build services using docker-compose
	```sh
	make start
	```

	This will start the services specified in `docker-compose.yml`.
	The values of the environment variable are set in `.makerc`.
	Typically this file wouldn't be commited but leaving it here for showcasing purposes

	There are 3 services:
	- answers (the actual project)
	- db (mongoDB database)
	- db-ui (mongoDB UI app, available at http://localhost:8081/)
		- the credentials are:
			- username: `me`
			- password: `ok`

- import [answers postman collection](./assets/answers.postman_collection) in `assets` folder
	- You will have to add an environment with the `base_url` environment variable to `http://localhost:8080` to make it work the first time.

- run `Scenarios` folder in the collection

- tear down the services
	```
	make down
	```


## Local development

Setup the typical Go environment. In this project I opted to use a VSCode [devcontainer](https://github.com/microsoft/vscode-dev-containers/tree/v0.245.0/containers/go).

Run the project by doing:
```
make run
```

This won't work by default without the db service being up, so ensure that is running first.

Alternatively, if you prefer to not use a database you can use the InMemory implementation by setting the environment variable `USE_IN_MEMORY` to `true` in the `.makerc` file. Bear in mind this implementation won't persist the data (I did it to test the API behaviour first).

Made a couple very simple unit tests with a generated mock using [genmock](https://pkg.go.dev/gitlab.com/so_literate/genmock).

To run the unit tests simply do:
```sh
make test
```

You can check code coverage by doing:
```sh
make test-coverage
```

This will generate the coverage.out and coverage.html files.
These files are ignore on the `.gitignore` file.

You can also generate the mocks by doing:
```
make mocks
```

## Questions

### **How would you support multiple users?**

Simplest way and assuming this would be a standalone app would be to add some basic auth. Some additional endpoints to create users with respective persistence layer with the usual security concerns (example: store hashed passwords and not in plain text).

Then when using the answers API group I'd add a middleware function to check and validate if the is user authorised to use the endpoint and access the data.

Every database query would have to take the userId into account to ensure the user doesn't get access to other users data.

### **How would you support answers with types other than string?**

I would switch the `value` property to an `interface{}` or potentially to named interface such as ` type IAnswer interface{}` if some common behaviours can be found (example: they all need to compute a value by doing function `Value() string`).

After unmarshalling the request payload I would then proceed to do a switch with a type assertion to map the payload to the proper struct.

### **What are the main bottlenecks of your solution?**

The codebase currently only supports using one database with 2 specified collections. If multiple user support is needed this would have to change.

There is only one node for the database so the more requests it takes the slower it will take to response over time.

It attempts to do database requests every single time an endpoint is hit. This can be vastly improved with the introduction of a caching layer.

### **How would you scale the service to cope with thousands of of requests?**

In this project I dockerized the application so be made into multiple replicas (example: increase the number of pods in a kubernetes cluster).

Instead of accessing one single node for the database it should access a cluster instead where ideally it would sharing the data via multiple servers (example: sharding)

Implement a shared caching layer (example: redis or memcached servers) would also greatly help with processing the requests.

---
## Task
Create a service that exposes an API (graphQL or REST) which allows people to create, update, delete and retrieve answers as key-value pairs. The answers should be stored in a persistent storage medium so they can handle service restarts.
An answer can be defined as:

```json
{
	"key" : "name",
	"value" : "John"
}
```

The API should expose the following endpoints:

- create answer
- update answer
- get answer (returns the latest answer for the given key)
- delete answer
- get history for given key (returns an array of events in chronological order)

An event can be defined as:

```json
{
	"event" : "create",
	"data" : {
		"key": "name",
		"value": "John"
	}
}
```

If a user saves the same key multiple times (using update), every answer should be saved. When retrieving an answer, it should return the latest answer.

If a user tries to create an answer that already exists - the request should fail and an adequate message or code should be returned.

If an answer doesn't exist or has been deleted, an adequate message or code should be returned.

When returning history, only mutating events (create, update, delete) should be returned. The "get" events should not be recorded.

It is possible to create a key after it has been deleted. However, it is not possible to update a deleted key. For example the following event sequences are allowed:

create → delete → create → update

create → update → delete → create → update

However, the following should not be allowed:

create → delete → update
create → create


## Additional questions:

- How would you support multiple users?
- How would you support answers with types other than string
- What are the main bottlenecks of your solution?
- How would you scale the service to cope with thousands of requests?

# How to submit:

- Send us a link to the source code repository.
- Make sure to include instructions on how to run the project in the `README.md`
- Include unit-tests and/or integration-tests (e.g. using Postman)