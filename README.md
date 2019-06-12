# Payments Sample [![Build Status](https://travis-ci.org/michaljemala/payments-sample.svg?branch=master)](https://travis-ci.org/michaljemala/payments-sample)[![Go Report Card](https://goreportcard.com/badge/github.com/michaljemala/payments-sample)](https://goreportcard.com/report/github.com/michaljemala/payments-sample)

This sample demonstrates a simple Payments server with a REST API exposed over HTTP.  

## API

API is modeled with REST principles in mind. A single resource to manage payments is exposed at the moment. By default payment instances are persisted using PostgreSQL, but this can be easily switched to any other RDBMS (or NoSQL database if desired). The API is described using OpenAPI v3, see the [specs file](./api/openapi.yaml) directly or run the server and navigate to `http://localhost:8080/docs`. It follows the [Zalando guidelines](https://opensource.zalando.com/restful-api-guidelines/) for defining RESTful APIs.   

### GET /payments
Retrieve collection of payments.

### GET /payments/{payment_id}
Retrieve an existing payment.

### POST /payments
Create a new payment.

### PATCH /payments/{payment_id}
Edit an existing payment.

### DELETE /payments/{payment_id}
Delete an existing payment.

## Run server 

You have two options, either use Docker Compose and run `docker-compose up` or run server locally `go run cmd/payments-server/main.go -http :8080 -database postgres:///payments -migrations file://./scripts/migrations/postgres`. In order to run server locally you have to have a running Postgres database server with a database named `payments` created. Server can be gracefully shut down by sending it the `SIGINT` or `SIGTERM` signals (just use `CTRL+C` when running locally).  

## Run tests
Codebase is unit-tested and dependencies are mocked so no database is required to be prepared, just run `go test ./...`.

## Project Structure
The acknowledged [standard Go project structure](https://github.com/golang-standards/project-layout) is used to avoid confusion.

## Payments Server Design
Exposed payment resource follows the [json:api](https://jsonapi.org) specification. The structure is reflected in the code: a [Payments API](./pkg/payments/api.go) defines a [single payment resource](./pkg/payments/resource.go) which ensures a correct payloads are being exchanged as per the json:api specification. The resource then delegates to a [payments service](./pkg/payments/service.go) which encapsulates business logic, validation and orchestrates the calls to stores (repositories). [Stores](./pkg/payments/store.go) ensures the resource's parts are correctly persisted and loaded to/from the database.

Each layer is abstracted using Go interfaces to allow easy unit-testing and enable transparently add/replace specific implementations.

As per Go best practices interfaces are implied, i.e. not being defined on the implementing side, but on the consuming side instead.

The Payments domain is described in the [domain package](./pkg/domain), which the defines primitive and complex objects used by the resource, service and stores.

All errors are being wrapped in to a single [common structure](./pkg/internal/errors) to allow proper conversions to HTTP statuses and json:api error payloads.
