package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChirpNewUserRequest struct {
	Email string
}

func apiNewUserHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := ChirpNewUserRequest{}
	err := decoder.Decode(&jsonReq)
	if err == nil {
		if len(jsonReq.Email) < 3 {
			respondWithError(w, 400, "email is too short")
			return
		}
		user, err := mainGLOBS.dbQueries.CreateUser(r.Context(), jsonReq.Email)
		if err != nil {
			fmt.Printf("new user db fail %v \n", err)
			respondWithError(w, 501, "internal server error")
			return
		}
		respondWithJSON(w, 201, user)
		return
	}
	fmt.Printf("apiNewUserHand %v", err)
}
