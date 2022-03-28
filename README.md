# pulse id service

The requirement is given in ./Platform Backend Engineer Task.pdf file

## Application structure and details

- Default environment variable is in .env file (which is checked-in just for the reference)
- Sqlite3 is used as a database  
- Default Admin user and table schema is inserted as a part of db script present in ./scripts/db_schema.sql  
- User need to be logged-in to access the private routes i.e.

1) api/v1/genToken
2) api/v1/getAllToken
3) api/v1/invalidateToken

- Below are the public routes and both of them are rate limited with 5 request at a time, after that user has to wait for 60Sec and for full reset user has to wait for 1 hour to make another request, currently this values are hardcoded in the code which can be made configurable

1) /login
2) /validatetoken

- App Token EXPIRE_IN_DAYS is passed though env file
- All the test cases for the repository and service are present in respective folder
- Docker and build pipeline is added

## How to run the application

### Locally 

To run the application do

```sh 
$ go mod download  
$ go run main.go
```

### Using docker

```sh 
$ docker build . -t pulseid 
$ docker run -p 8081:8081 -it pulseid 
```

### Accessing the routes

1) User tried to generate token without login

``` 
request :-> curl --location --request GET 'http://localhost:8081/api/v1/genToken'
response:-> token not present
```

2) User has logged-in with the username and password inserted in default db script

```
request :->  curl --location --request POST 'http://localhost:8081/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username" : "admin",
    "password" : "admin"
}'

response:->  {
    "role": "ADMIN",
    "username": "admin",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Mjk2MDQyMjQsInJvbGUiOiJBRE1JTiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.bSSPrJF1bubv2IcMsMSfE7S4_-TUVGy8i8EkT_cQ15A"
}
```

3) User the above token to access the private routes, To generate the token

```
request :-> curl --location --request GET 'http://localhost:8081/api/v1/genToken' \
--header 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Mjk2MDQyMjQsInJvbGUiOiJBRE1JTiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.bSSPrJF1bubv2IcMsMSfE7S4_-TUVGy8i8EkT_cQ15A'


response:-> "PpPRm3p9GyN"
```

4) Validate the app token

``` 
request :->  curl --location --request GET 'http://localhost:8081/validatetoken/GDVTL8ipbZS'
response:->   true
```

5) Check if error message is sent when invalid app token is passed

```
request :-> curl --location --request GET 'http://localhost:8081/validatetoken/sdgfasdg'
response:->  record not found , 400

```

6) Check if length of app token is incorrect

```
request :->  curl --location --request GET 'http://localhost:8081/validatetoken/12345'
response:-> invalid app token, 400
```

7) Get all token Active and Inactive

```
request :-> curl --location --request GET 'http://localhost:8081/api/v1/getAllToken' \
--header 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Mjk2MzEyMzgsInJvbGUiOiJBRE1JTiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LZRjU6W_gdbscmjfNfzWuqecYpvCDPxWV2nnzJpZqBs'

response:->  [
    {
        "id": "af0ba00b-570b-47a0-91b3-41766bd3e16a",
        "username": "admin",
        "token": "PpPRm3p9GyN",
        "exp_date": "2021-08-29T09:06:30.664079087+05:30",
        "is_active": true
    },
    {
        "id": "5fa9b709-e230-4df4-82b4-313cdddcc7ce",
        "username": "admin",
        "token": "HAQYVEwhbsq",
        "exp_date": "2021-08-22T09:37:47.114203267+05:30",
        "is_active": false
    },
    {
        "id": "051dd1b7-a3d5-4427-8910-09c3be45296c",
        "username": "admin",
        "token": "ZKww2WFHQfL",
        "exp_date": "2021-08-23T09:38:33.435709572+05:30",
        "is_active": true
    }
]
```

8) Deactivate the token

```
request :->  curl --location --request PATCH 'http://localhost:8081/api/v1/invalidateToken' \
--header 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Mjk2MzEyMzgsInJvbGUiOiJBRE1JTiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LZRjU6W_gdbscmjfNfzWuqecYpvCDPxWV2nnzJpZqBs' \
--header 'Content-Type: application/json' \
--data-raw '{
    "appToken" : "pbSKYdfHllo"
}'

response :-> {
    "id": "0a9b454b-df74-436f-a809-0b6cd6203408",
    "username": "admin",
    "token": "pbSKYdfHllo",
    "exp_date": "2021-08-23T16:24:57.350421438+05:30",
    "is_active": false
}
```
