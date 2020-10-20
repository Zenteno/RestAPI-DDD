# API Rest Example

Ultra fast api rest to handle history sessions, built with GoFiber.

## Getting Started

These instructions will get you a copy of the project up and running on any computer with Docker and Docker-compose installed.

### Prerequisites

For development you will need:

```
Go 1.15+
Mongodb
```

Or if you just want to run and deploy:

```
Docker
Docker-Compose
Make
```

### Installing

Of course the first thing to do is to clone repository

```
git clone https://github.com/Zenteno/RestAPI-DDD
```

And to start our software
```
make up
```
After that, everything is good to go and run locally and open your browser in http://localhost:3000/

## Running the tests

By default the test are executed using the same mongo database, but with a different collection name.

To run the test:

```
make test
```
This command will generate a coverage html file in the project folder, named "coverage.html", where you can see all covered lines with the test file.

![](./assets/coverage.png?raw=true)

## API Paths

### NewLocation
**You send:**  A new Session History.

**Request:**
```json
POST /api/v1/location HTTP/1.1
Accept: application/json
Content-Type: application/json
Content-Length: xy

{
	"shopper_uuid":  "b6549bd2-0c33-11eb-adc1-0242ac120002",
	"session_uuid": "beef2eb0-0c33-11eb-adc1-0242ac120002" ,
	"lat": -36.8126209 ,
	"lng": -73.0392838,
	"precision":  12.3,
	"reported_at": "2020-10-12T02:35:57+00:00",
}
```
**Successful Response:**
```json
HTTP/1.1 200 OK
Server: My RESTful API
Content-Type: application/json
Content-Length: xy

```
**Failed Response:**
```json
HTTP/1.1 400 Bad Request
Server: My RESTful API
Content-Type: application/json
Content-Length: xy

{
    "message": "missing parameters"
}
``` 
### History Session
**You send:**  A session uuid.

**Request:**
```json
POST /api/v1/session_location_history/{session_uuid} HTTP/1.1
Accept: application/json
Content-Length: xy

```
**Successful Response:**
```json
HTTP/1.1 200 OK
Server: My RESTful API
Content-Type: application/json
Content-Length: xy

[
	{
		"shopper_uuid":  "b6549bd2-0c33-11eb-adc1-0242ac120002",
		"session_uuid": "beef2eb0-0c33-11eb-adc1-0242ac120002" ,
		"lat": -36.8126209 ,
		"lng": -73.0392838,
		"precision":  12.3,
		"reported_at": "2020-10-12T02:35:57+00:00",
	}
]

```
**Failed Response:**
```json
HTTP/1.1 404 Not Found
Server: My RESTful API
Content-Type: application/json
Content-Length: xy
``` 

### Current Location
**You send:**  A shopper uuid.

**Request:**
```json
POST /api/v1/current_shopper_location/{shopper_uuid} HTTP/1.1
Accept: application/json
Content-Length: xy

```
**Successful Response:**
```json
HTTP/1.1 200 OK
Server: My RESTful API
Content-Type: application/json
Content-Length: xy

{
		"shopper_uuid":  "b6549bd2-0c33-11eb-adc1-0242ac120002",
		"session_uuid": "beef2eb0-0c33-11eb-adc1-0242ac120002" ,
		"lat": -36.8126209 ,
		"lng": -73.0392838,
		"precision":  12.3,
		"reported_at": "2020-10-12T02:35:57+00:00",
}


```
**Failed Response:**
```json
HTTP/1.1 404 Not Found
Server: My RESTful API
Content-Type: application/json
Content-Length: xy
``` 

## Swagger
The API also comes with a swagger integration, you can find it in http://localhost:3000/docs

![](./assets/swagger.png?raw=true)

## Built With

* [Fiber](https://gofiber.io/) 
* [GoFakeit](https://github.com/brianvoe/gofakeit)
* [Logrus](https://github.com/sirupsen/logrus)
* [Testify](https://github.com/stretchr/testify)
* [Swaggo](https://github.com/swaggo/swag )
* [Validator](https://github.com/go-playground/validator/v10 )

## Author

* **Alberto Zenteno**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


