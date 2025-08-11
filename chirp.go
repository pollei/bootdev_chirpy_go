package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	//"github.com/pollei/bootdev_chirpy_go/internal/database"
	"github.com/pollei/bootdev_chirpy_go/internal/auth"
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
	if err != nil {
		fmt.Printf("new chirp decode fail %v \n", err)
		respondWithError(w, 501, "internal server error")
		return
	}

	if len(jsonReq.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	bearTok, err := auth.GetBearerToken(r.Header)
	if err == nil && len(bearTok) > 5 {
		uuid, err := auth.ValidateJWT(bearTok, mainGLOBS.jwtSecretKey)
		if err == nil {
			jsonReq.UserID = uuid
		} else {
			respondWithError(w, 401, "jwt not valid")
			return
		}
	}
	clean := cleanPhrase(jsonReq.Body)
	jsonReq.Body = clean
	chirp, err := mainGLOBS.dbQueries.CreateChirp(r.Context(), jsonReq)
	if err != nil {
		fmt.Printf("new chirp db fail \n\t %v \n\t%v \n", jsonReq, err)
		respondWithError(w, 501, "internal server error")
		return
	}
	respondWithJSON(w, 201, chirp)
}

func apiGetAllChirpsHand(w http.ResponseWriter, r *http.Request) {
	chirps, err := mainGLOBS.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		fmt.Printf("apiGetAllChirpsHand db fail %v \n", err)
		respondWithError(w, 501, "internal server error")
		return
	}
	respondWithJSON(w, 200, chirps)
}

func apiGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirp_id")
	fmt.Printf(" apiGetChirpById %s \n", chirp_id)
	chirp_uuid, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 404, "not valid chirp id")
		return
	}
	fmt.Printf(" apiGetChirpById %s valid format \n", chirp_id)
	chirp, err := mainGLOBS.dbQueries.GetChirpByID(r.Context(), chirp_uuid)
	if err == sql.ErrNoRows {
		respondWithError(w, 404, "not valid chirp id")
		return
	}
	fmt.Printf(" apiGetChirpById db err %v \n", err)
	if err == nil {
		fmt.Printf(" apiGetChirpById sending %v", chirp)
		respondWithJSON(w, 200, chirp)
		return
	}
	//noRowErr := sql.noRowErr
}

func apiDeleteChirpByIdHand(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirp_id")
	fmt.Printf(" apiDeleteChirpById %s \n", chirp_id)
	bearTok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithEmpty(w, 401)
		return
	}
	bearUuid, err := auth.ValidateJWT(bearTok, mainGLOBS.jwtSecretKey)
	if err != nil {
		respondWithEmpty(w, 401)
		return
	}
	chirp_uuid, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 404, "not valid chirp id")
		return
	}
	fmt.Printf(" apiDeleteChirpById %s valid format \n", chirp_id)
	chirp, err := mainGLOBS.dbQueries.GetChirpByID(r.Context(), chirp_uuid)
	if err == sql.ErrNoRows {
		respondWithError(w, 404, "not valid chirp id")
		return
	}
	if chirp.UserID != bearUuid {
		respondWithEmpty(w, 403)
		return
	}
	delParams := database.DeleteOwnChirpByIDParams{
		ID: chirp_uuid, UserID: bearUuid}
	_, err = mainGLOBS.dbQueries.DeleteOwnChirpByID(r.Context(), delParams)
	if err == sql.ErrNoRows {
		respondWithEmpty(w, 404)
		return
	}
	respondWithEmpty(w, 204)
	//fmt.Printf(" apiDeleteChirpById db err %v \n", err)
}

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
