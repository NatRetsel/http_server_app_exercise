package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/natretsel/http_server_app_exercise/internal/auth"
	"github.com/natretsel/http_server_app_exercise/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	/*
		Expects login credentials
		{
			"email": "${user email}",
			"password": "${password}"
		}

		1.) Unmarshal JSON
		2.) Check if account exist with email in DB
		3.) Check if password hash matches
		4.) Create JWT (2h) and refresh token (60 days)
		5.) Send tokens back as response
		{
			User: user struct without hashed_pw
			Token: access token
			RefreshToken: refresh token
		}
	*/
	type LoginParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Unmarshal JSON
	loginParam := LoginParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not unmarshal JSON", err)
		return
	}

	// Query for user
	user, err := cfg.db.GetUserByEmail(r.Context(), loginParam.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user does not exist", err)
		return
	}

	// Check hashed PW
	err = auth.CheckPasswordHash(user.HashedPassword, loginParam.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	// create JWT
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*2)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create access token", err)
		return
	}

	// create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create refresh token", err)
		return
	}

	// Store refresh token in DB
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
		return
	}

	type LoginResponse struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	// Generate response
	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
