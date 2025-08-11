package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	//"time"

	"github.com/google/uuid"
	"github.com/pollei/bootdev_chirpy_go/internal/auth"
	//"github.com/pollei/bootdev_chirpy_go/internal/database"
)

type PolkaApiWebhookData struct {
	UserID uuid.UUID `json:"user_id"`
}

type PolkaApiWebhook struct {
	Event string              `json:"event"`
	Data  PolkaApiWebhookData `json:"data"`
}

func apiPolkaWebhooksHand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonReq := PolkaApiWebhook{}
	err := decoder.Decode(&jsonReq)
	if err != nil {
		respondWithEmpty(w, 404)
		return
	}
	fmt.Printf(" apiPolkaWebhooksHand %v \n", jsonReq)
	polkaKey, err := auth.GetApikey(r.Header)
	if err != nil || polkaKey != mainGLOBS.polkaWebhookKey {
		respondWithEmpty(w, 401)
		return
	}
	if jsonReq.Event != "user.upgraded" {
		respondWithEmpty(w, 204)
		return
	}
	_, err = mainGLOBS.dbQueries.UpdateUserToRedByID(
		r.Context(), jsonReq.Data.UserID)
	if err == sql.ErrNoRows {
		respondWithEmpty(w, 404)
		return
	}
	respondWithEmpty(w, 204)

}
