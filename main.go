package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pollei/bootdev_chirpy_go/internal/database"
)

type mainEvilGlobals struct {
	db           *sql.DB
	dbQueries    *database.Queries
	platform     string
	jwtSecretKey string
	hcHand       *handlerHitCounter
}

var mainGLOBS mainEvilGlobals

func adminResetHand(w http.ResponseWriter, r *http.Request) {
	if mainGLOBS.platform != "dev" {
		respondWithError(w, 403, "")
		return
	}
	fmt.Printf("adminResetHand plat %s \n", mainGLOBS.platform)
	mainGLOBS.dbQueries.DeleteAllUsers(r.Context())
	mainGLOBS.hcHand.Reset()
	respondWithEmpty(w, 200)
}

func main() {
	srvMux := http.NewServeMux()
	hcHand := newHandlerHitCounter(http.StripPrefix(
		"/app/",
		http.FileServer(http.Dir("public_html"))))
	mainGLOBS.hcHand = hcHand
	srvMux.Handle("/app/", hcHand)
	srvMux.HandleFunc("POST /admin/reset", adminResetHand)
	//srvMux.HandleFunc("POST /api/validate_chirp", validateCleanChirpHand)
	srvMux.HandleFunc("POST /api/users", apiNewUserHand)
	srvMux.HandleFunc("PUT /api/users", apiUserHand)
	srvMux.HandleFunc("POST /api/login", apiLoginUserHand)
	srvMux.HandleFunc("POST /api/refresh", apiRefreshHand)
	srvMux.HandleFunc("POST /api/revoke", apiRevokeHand)
	srvMux.HandleFunc("POST /api/chirps", apiNewChirpHand)
	srvMux.HandleFunc("DELETE /api/chirps/{chirp_id}", apiDeleteChirpByIdHand)
	srvMux.HandleFunc("GET /api/chirps", apiGetAllChirpsHand)
	srvMux.HandleFunc("GET /api/chirps/{chirp_id}", apiGetChirpById)
	srvMux.Handle("POST /admin/metrics_reset", hcHand.reset)
	srvMux.Handle("GET /admin/metrics", hcHand.metrics)
	srvMux.HandleFunc("GET /api/healthz", healthHand)
	httpD := &http.Server{
		Handler: srvMux,
		Addr:    ":8080",
	}
	dbURL := os.Getenv("DB_URL")
	mainGLOBS.platform = os.Getenv("PLATFORM")
	mainGLOBS.jwtSecretKey = os.Getenv("JWT_SECRET")
	if len(mainGLOBS.platform) < 1 {
		mainGLOBS.platform = "prod"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		os.Exit(1)
	}
	mainGLOBS.db = db
	mainGLOBS.dbQueries = database.New(db)
	httpD.ListenAndServe()
}
