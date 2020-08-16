# Go-API

A simple API service built with Go.

**Installing**

Go installation guide can be found [here](https://golang.org/doc/install)

```sh
# Install project dependencies
go get -u ./...
```

**Running**

```sh
go run .\main.go .\server.go
```

## Usage

All responses will have the same form.

```json
{
    "status": "Integer holding the status code of the response",
    "state": "String that will either be `ok` or `fail`",
    "result": "Mixed type holding the content of the response",
}
```

Responses definations will only show the value of the `result field`*.

### Retriving a Screenshot

**Definition**

`GET /screenshots/{identifier}`

**Arguments**

- `"key":string` the key which will be validated
- `"name":string` the name of the image

**Response**

- `200 OK` on success
- `404 Not found` if screenshot could not be found
- `401 Unauthourzed` if key is invalid

```json
{
	"status": 200,
	"state": "ok",
	"result": {
		"img": "Binary Image",
		"name": "Amazing-Bird-200",
		"mime": "image/jpg",
		"timestamp": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
	}
}
```
```json
{
	"status": 404,
	"state": "fail",
	"result": "error: could not find screenshot with name 'Ayyy'"
}
```
```json
{
	"status": 401,
	"state": "fail",
	"result": "error: key is unauthorized"
}
```

### Creating a Screenshot

**Definition**

`POST /screenshots`

**Arguments**

- `"key":string` the key which will be validated
- `"img":multipart/form-data` the image which will be saved
- `"name":string` the name of the image
- `"mime":string` the mime type of the image
- `"timestamp":string (optional)` the tiemstamp, will be generated if not provided


**Response**

- `201 Created` on success
- `409 Conflict` on duplicate names
- `401 Unauthourzed` if key is invalid

```json
{
	"status": 409,
	"state": "fail",
	"result": "error: screenshot with name 'cool name' already exists"
}
```

### Removing a Screenshot

**Definition**

- `DELETE /screenshots`

**Response**

- `200 OK` on success
- `404 Not Found` if screenshot could not be found 
- `401 Unauthourzed` if key is invalid

```json
{
	"status": 200,
	"state": "ok",
	"result": {
		"img": "Binary Image",
		"name": "Amazing-Bird-200",
		"mime": "image/jpg",
		"timestamp": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
	}
}
```