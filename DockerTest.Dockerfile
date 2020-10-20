# Building the binary of the App
FROM golang:1.15 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /app

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

RUN go get -u github.com/swaggo/swag/cmd/swag

RUN swag init -g app.go

RUN go test -coverpkg="api-ddd/..." -c -p 1 -tags testrunmain api-ddd

RUN mkdir -p /app/coverage