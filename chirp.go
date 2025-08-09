package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	//"github.com/pollei/bootdev_chirpy_go/internal/database"
	"github.com/pollei/bootdev_chirpy_go/internal/database"
)

type ChirpPostRequest struct {
	Body   string
	UserId uuid.UUID
}
type ChirpPostResponse struct {
	Error       string `json:"error,omitempty"`
	Valid       bool   `json:"valid,omitzero"`
	CleanedBody string `json:"cleaned_body,omitempty"`
}

type wordSet map[string]struct{}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp := ChirpPostResponse{}
	jsonResp.Error = msg
	w.WriteHeader(code)
	dat, _ := json.Marshal(jsonResp)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, _ := json.Marshal(payload)
	w.Write(dat)

}
func respondWithEmpty(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func cleanPhrase(in string) string {
	badWords := wordSet{
		"kerfuffle": struct{}{},
		"sharbert":  struct{}{},
		"fornax":    struct{}{},
	}
	arr := strings.Split(in, " ")
	for i, word := range arr {
		if _, isBad := badWords[strings.ToLower(word)]; isBad {
			arr[i] = "****"
		}
	}
	return strings.Join(arr, " ")
}

func apiNewChirpHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := database.CreateChirpParams{}
	err := decoder.Decode(&jsonReq)
	if err == nil {
		if len(jsonReq.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		clean := cleanPhrase(jsonReq.Body)
		jsonReq.Body = clean
		chirp, err := mainGLOBS.dbQueries.CreateChirp(r.Context(), jsonReq)
		if err != nil {
			fmt.Printf("new user db fail %v \n", err)
			respondWithError(w, 501, "internal server error")
			return
		}
		respondWithJSON(w, 201, chirp)
		return
	}
	fmt.Printf("apiNewChirpHand exit %v", err)
}

//nolint:unused
func validateCleanChirpHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpPostRequest{}
	err := decoder.Decode(&jsonReq)
	if err == nil {
		if len(jsonReq.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		clean := cleanPhrase(jsonReq.Body)
		jsonResp := ChirpPostResponse{Valid: true, CleanedBody: clean}
		respondWithJSON(w, 200, jsonResp)
	}

}

func validateChirpHand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpPostRequest{}
	jsonResp := ChirpPostResponse{}
	err := decoder.Decode(&jsonReq)
	if err == nil {
		if len(jsonReq.Body) > 140 {
			jsonResp.Error = "Chirp is too long"
			w.WriteHeader(400)
			dat, _ := json.Marshal(jsonResp)
			w.Write(dat)
			return
		}
		jsonResp.Valid = true
		dat, _ := json.Marshal(jsonResp)
		w.WriteHeader(200)
		w.Write(dat)
		return
	}
}
