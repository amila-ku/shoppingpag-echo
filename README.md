# shoppingpal

CRUD API using golang and dynamodb as backend database.

## status

![](https://github.com/amila-ku/shoppingpal/workflows/Go/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/amila-ku/shoppingpal)](https://goreportcard.com/report/github.com/amila-ku/shoppingpal)


## Local Setup


start dynamodb locally

```
docker run -d -p 8000:8000 amazon/dynamodb-local

```

using aws cli to check tables

```
aws dynamodb list-tables --endpoint-url http://localhost:8000

```

using cli to list table content

```
aws dynamodb scan --table-name itemtable  --endpoint-url http://localhost:8000

```
