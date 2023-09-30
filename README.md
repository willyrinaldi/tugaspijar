# Simple Golang API

## Documents

Echo framework https://echo.labstack.com/guide

Go driver for MongoDB https://github.com/mongodb/mongo-go-driver

## MongoDB container

Create volume

```bash
docker volume create simple_mongodbdata
```

Start monodb container

```bash
docker run -it --name simple-mongo -d -p 27017:27017 -v simple_mongodbdata:/data/db mongo
atau
docker-compose up

```

## Run simple API

```bash
go run main.go
```

## Endpoint

Browse to http://localhost:1323