package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/natretsel/http_server_app_exercise/internal/auth"
	"github.com/natretsel/http_server_app_exercise/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	/*
		Expect as JSON input
		{
			"email" : "${user email}"
			"password": "${user password}"
		}
		1.) Unmarshal the JSON input
		2.) Query for existing user with the same email
			- return bad request if user with email already exist
		3.) Hash the password
			- return bad request if we could not hash password
		4.) Write the user account details into the DB

		- WIP: If authenticated user access this endpoint, redirect to home page
	*/
	type UserParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userParams := UserParameters{}
	decoder := json.NewDecoder(r.Body) // read from request body
	err := decoder.Decode(&userParams) // unpack into userParams
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error unmarshalling JSON", err)
		return
	}
	// Check for existing user with the provided email
	user, err := cfg.db.GetUserByEmail(r.Context(), userParams.Email)
	if err == nil || user.Email == userParams.Email {
		respondWithError(w, http.StatusBadRequest, "user with email already exist", errors.New("user with email already exist"))
		return
	}
	// Check for empty fields
	if userParams.Password == "" || userParams.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}
	// Check for appropriate format - WIP - endpoint accessible from API calls

	// Hash password
	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password, try a different password", err)
		return
	}
	// Add records in DB
	user, err = cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userParams.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user record in DB", err)
		return
	}
	// return OK
	type UserResponse struct {
		User `json:"user"`
	}
	respondWithJSON(w, http.StatusCreated, UserResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
