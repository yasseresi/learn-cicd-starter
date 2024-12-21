package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		respondWithError(w, http.StatusBadRequest, "Empty request body")
		return
	}

	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	apiKey, err := generateRandomSHA256Hash()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't gen apikey")
		return
	}

	err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Name:      params.Name,
		ApiKey:    apiKey,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}
	respondWithJSON(w, http.StatusCreated, userResp)
}

func generateRandomSHA256Hash() (string, error) {
	const keyLength = 32
	randomBytes := make([]byte, keyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:]), nil
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
