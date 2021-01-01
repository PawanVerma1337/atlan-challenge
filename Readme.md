# Atlan Backend Task
A prototype code displaying the resumable uploads using the tus protocol + a terminate upload route. The upload work by appending a binary file to the patch request. The post request creates a file in the database. The subsequent patch request is used to upload data to the server and if the connection is dropped an offset is saved to database in order to resume the upload from the offset.

## Building from Dockerfile

```
docker build -t atlan-challenge .
```
To run
```
docker run atlan-challenge
```

## Routes

```
Routes                      |   Method
----------------------------|-----------------
/upload                     |   (POST)
/upload/resume?id={id}      |   (PATCH)
/upload/terminate?id={id}   |   (GET)
/upload/status?id={id}      |   (GET)
```

## Stack
- Golang
- Mongo
- Gorilla/Mux (Router)
