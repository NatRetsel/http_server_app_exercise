# http_server_app_exercise
Thought process in creating a http server app

This project aims to build the thought process behind a Web app requiring a http server so that anyone can build it in any language - know what is required and look for the components!

Let's use an example to illustrate - a youtube clone in Golang. We'll call it ViewTube. Let's keep it simple, the app shall:
* Allow new users to register an account
* Allow users to upload videos, along with thumbnails, title and description.
* Allow users to delete their own uploaded videos

# Thought process

## Backend
I will start with the backend. This section touches on the API and endpoints which the frontend will consume, as well as the database and authentication logic.

### Server
The server will need to cater for various API and endpoints:

| Requirement                                                   | HTTP Method | Endpoint                                                    | Purpose                                                      |
| ------------------------------------------------------------- | ----------- | ----------------------------------------------------------- | ------------------------------------------------------------ |
| Access the application                                        | GET         | `/app`                                                      | Arrive at homepage of ViewTube                               |
| Allow new users to register an account                        | POST        | `/api/users`                                                | Register new user account, store information in DB.          |
| Allow registered users to login                               | POST        | `/api/login`                                                | Authenticate registered users.                               |
| Allow users to upload videos, thumbnails, title, description. | POST        | `/api/videos/{videoID}`<br>`/api/thumnail_upload/{videoID}` | Auth endpoint allowing users to upload videos and thumbnails |
| Allow users to delete their own uploaded videos               | DELETE      | `/api/videos/{videoID}`                                     | Auth endpoint allowing users to delete their uploaded video  |
| Allow users to retrieve their uploaded videos                 | GET         | `/api/videos`                                               | Users should see their list of uploaded videos               |
While our server listens to request on our specified port, we need handlers for each of the endpoints. 

We can choose:
* To use the default handler, which implicitly uses `http.DefaultServeMux`
* Or instantiate with explicit `http.DefaultServeMux` and `http.Server` struct for control, isolation and custom server config. `http.Server` struct is optional.
We will go with the latter.

```go
package main

import (
	"net/http"
	"fmt"
)

type apiConfig struct {
	// config variables here, such as DB, s3Client, jwt, directories
}


func main(){
	// read env variables ...

	port := "8080"
	// instantiate apiConfig struct
	cfg := apiConfig{
		// env variables
	}
	// instantiate mux and http server
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v",port),
		Handler: mux,
	}
	// add route handlers here
	mux.HandleFunc("GET /app", http.FileServer(http.Dir(filePathRoot)))

	mux.HandleFunc("POST /api/users", cfg.someFunctionImplementingHandlerServeHTTP)


	srv.ListenAndServe()
}
```

We will serve our app from the `./app` directory (`filePathRoot`). The frontend components like submit buttons will trigger the request to our API endpoints, supplying information captured in form as JSON. We will also specify information in the request header, such as bearer tokens.

The purpose of the `apiConfig` struct is a way to capture the application state variables. We can then simply implement handlers as methods to this struct. We could also do without the struct and wrap the handlers with a middleware and pass in the additional required variables.

### Database
We want to match users with their uploaded videos - therein a relation between users and their uploaded videos. We will use a relational database for this project - Postgresql. For auth purposes, we will also have a `refresh_token` table.

![[ViewTube.png]]

To save time in writing type safe queries, we will use `sqlc` command line tool. It requires a `sqlc.yaml` file to know where to read the schema, queries and where to output the code. We will include this `sqlc.yaml` file in the root directory of our project.

We will also need a database migration script - `goose` to facilitate ongoing development and future changes to the database schema.

### Authentication

#### Access token
An access token is a credential that a client uses to access protected resources on a server. This can come in the form of a JSON Web Token (JWT) issued by an authorization server after a user successfully logs in or authenticates.

The token contains information like the user's identity and permissions and is sent with each request to prove the client is allowed to access the requested resource. Access tokens usually have a short lifespan (minutes to hours) for security reasons.

In Web Apps, it is typically stored in memory (e.g. a JS variable) or, less securely, in `localStorage` or `sessionStorage`. Memory is preferred because it's cleared when the tab closes, reducing exposure. It is sent to the server in the HTTP `Authorization` header as a Bearer token (`Authorization: Bearer <acces_token>`).

To create a JWT in Golang:

```go
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){

	signingKey := []byte(tokenSecret)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{

		Issuer: string(TokenTypeAccess),
	
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
	
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	
		Subject: userID.String(),
	})

	return token.SignedString(signingKey)
}
```

A JWT has three parts, **Base64 encoded** and separated by dots:
`header.payload.signature`

`header`:
* Contains the algorithm and token type `{"alg": "HS256", "typ":"JWT"}`

`payload`:
* Contains the claims. E.g. `{"subject":"${userId}", "exp":"${expire datetime} "iat":"${issue datetime}"`

`signature`:
* A hash:
	* The encoded header and payload are concatenated with a dot (`<base64_header>.<base64_payload>`) and fed into the specified cryptographic algorithm with the secret key.
	* This produces a fixed length hash serving as the signature.
	* The signature is then Base64URL-encoded and appended to the token (`<base64_header>.<base64_payload>.<base64_signature>`)

* We need a claims - essentially the payload of the JWT. It contains the information about the user.
* Specifying the signing method, we get a token
* Then we sign the token with our signing key (secret token in bytes)

If the payload is tempered, the signature won't match when recomputed with the secret key

To verify a JWT Token, we will:
* Ideally have the claims struct.
* The header has information on the algorithm used, so we will call `jwt.ParseWithClaims`, passing in the `tokenString` we want to validate and a function that returns the secret key used during signing. 
* The decoded payload JSON is unmarshaled into the claims struct.
* The validation of signature - the method recomputes the signature over `<header>.<payload>`
	* The recomputed hash is compared to the decoded `signature` from the token.

```go
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {


	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct,
	
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
	
	if err != nil {
		return uuid.Nil, err
	}
	
	userIDString, err := token.Claims.GetSubject()
	
	if err != nil {
		return uuid.Nil, err
	}
	
	issuer, err := token.Claims.GetIssuer()
	
	if err != nil {
		return uuid.Nil, err
	}
	
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	
	id, err := uuid.Parse(userIDString)
	
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil

}
```

#### Refresh token
Refresh token is a longer lived credential used to obtain a new access token once the original one expires, without requiring the user to log in again. It's issued alongside the access token during authentication and stored securely by the client.

When the access token expires, the client sends the refresh token to the authorization server, which verifies it and issues a new access token (and sometimes new refresh token). Refresh tokens typically last days, weeks, or even months, but they can be manually revoked if compromised.

In Web Apps, it is ideally stored in an HTTP-only, secure cookie. This prevents JS access and ensures it's only sent over HTTPS. `localStorage` or `sessionStorage` are riskier options. It is typically sent to the server in the HTTP request body (e.g. as a POST parameter) or as an HTTP-only cookie.

Since the refresh token does not hold information on the users, we'll encoding random bytes into hexadecimal string for our refresh token. Details on the user, expiry and revoked status will be stored in the DB.

```go
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
```

### Video upload