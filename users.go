package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/pollei/bootdev_chirpy_go/internal/auth"
	"github.com/pollei/bootdev_chirpy_go/internal/database"
)

type ChirpApiUserRequest struct {
	Email            string
	Password         string
	ExpiresInSeconds int
}

type ChirpApiUserResponse struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	Token          string    `json:"token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	HashedPassword string    `json:"-"`
}

type ChirpApiRefreshResponse struct {
	Token string `json:"token,omitempty"`
}

func apiNewUserHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpApiUserRequest{}
	err := decoder.Decode(&jsonReq)
	if err == nil {
		if len(jsonReq.Email) < 3 {
			respondWithError(w, 400, "email is too short")
			return
		}
		hsPw, err := auth.HashPassword(jsonReq.Password)
		if err != nil {
			respondWithEmpty(w, 404)
			return
		}
		dbReq := database.CreateUserParams{Email: jsonReq.Email, HashedPassword: hsPw}
		userDb, err := mainGLOBS.dbQueries.CreateUser(r.Context(), dbReq)
		if err != nil {
			fmt.Printf("new user db fail %v \n", err)
			respondWithError(w, 501, "internal server error")
			return
		}
		user := ChirpApiUserResponse{
			ID: userDb.ID, CreatedAt: userDb.CreatedAt,
			UpdatedAt: userDb.UpdatedAt, Email: userDb.Email}
		respondWithJSON(w, 201, user)
		return
	}
	fmt.Printf("apiNewUserHand %v", err)
}

func apiUserHand(w http.ResponseWriter, r *http.Request) {
	bearTok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "bearer token not found")
		//apiNewUserHand(w, r)
		return
	}
	if len(bearTok) < 5 {
		respondWithError(w, 401, "bearer token not found")
		return
	}
	uuid, err := auth.ValidateJWT(bearTok, mainGLOBS.jwtSecretKey)
	if err != nil {
		respondWithError(w, 401, "jwt not valid")
		return
	}
	oldUsrDb, err := mainGLOBS.dbQueries.GetUserByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 401, "uuid not valid")
		return
	}
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpApiUserRequest{}
	err = decoder.Decode(&jsonReq)
	if err != nil {
		respondWithEmpty(w, 401)
	}
	if len(jsonReq.Email) < 3 {
		jsonReq.Email = oldUsrDb.Email
	}
	hsPw, err := auth.HashPassword(jsonReq.Password)
	if err != nil {
		hsPw = oldUsrDb.HashedPassword
	}

	updParam := database.UpdateUserByIDParams{
		ID: uuid, Email: jsonReq.Email, HashedPassword: hsPw}
	userDb, err := mainGLOBS.dbQueries.UpdateUserByID(r.Context(), updParam)
	if err != nil {
		respondWithError(w, 501, "internal server error")
		return
	}
	user := ChirpApiUserResponse{
		ID: userDb.ID, CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt, Email: userDb.Email}
	respondWithJSON(w, 200, user)
}

func apiLoginUserHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpApiUserRequest{}
	err := decoder.Decode(&jsonReq)
	if err != nil {
		fmt.Printf("new login decode fail %v \n", err)
		respondWithError(w, 501, "internal server error")
		return
	}
	userDb, err := mainGLOBS.dbQueries.GetUserByEmail(r.Context(), jsonReq.Email)
	if err != nil {
		respondWithError(w, 401, "bad login")
		return
	}
	err = auth.CheckPasswordHash(jsonReq.Password, userDb.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "bad login")
		return
	}
	jwtTok, err := auth.MakeJWT(userDb.ID, mainGLOBS.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, 401, "bad login")
		return
	}
	rfshTok, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("new login make refresh token fail %v \n", err)
		respondWithError(w, 501, "internal server error")
		return
	}
	now := time.Now().UTC()

	rfrshParam := database.CreateRefreshTokenParams{
		Token: rfshTok, UserID: userDb.ID, ExpiresAt: now.Add(60 * 24 * time.Hour)}
	mainGLOBS.dbQueries.CreateRefreshToken(r.Context(), rfrshParam)
	user := ChirpApiUserResponse{
		ID: userDb.ID, CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt, Email: userDb.Email,
		Token: jwtTok, RefreshToken: rfshTok,
	}
	respondWithJSON(w, 200, user)

}

func apiRefreshHand(w http.ResponseWriter, r *http.Request) {
	bearTok, err := auth.GetBearerToken(r.Header)
	if err != nil && len(bearTok) < 5 {
		respondWithError(w, 401, "bearer token not found")
		return
	}
	rfshTokDb, err := mainGLOBS.dbQueries.GetRefreshTokenByToken(r.Context(), bearTok)
	if err != nil || rfshTokDb.RevokedAt.Valid {
		respondWithError(w, 401, "bearer token not valid")
		return
	}
	//now := time.Now().UTC()
	jwtTok, err := auth.MakeJWT(rfshTokDb.UserID, mainGLOBS.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, 501, "internal server error")
		return
	}
	resp := ChirpApiRefreshResponse{Token: jwtTok}
	respondWithJSON(w, 200, resp)

}

func apiRevokeHand(w http.ResponseWriter, r *http.Request) {
	bearTok, err := auth.GetBearerToken(r.Header)
	if err != nil && len(bearTok) < 5 {
		respondWithError(w, 401, "bearer token not found")
		return
	}
	mainGLOBS.dbQueries.RevokeRefreshToken(r.Context(), bearTok)
	respondWithEmpty(w, 204)
}
