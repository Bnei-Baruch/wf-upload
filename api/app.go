package api

import (
	"context"
	"encoding/json"
	"github.com/Bnei-Baruch/wf-upload/common"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type App struct {
	Router        *mux.Router
	tokenVerifier *oidc.IDTokenVerifier
}

func (a *App) InitClient() {
	oidcProvider, err := oidc.NewProvider(context.TODO(), common.ACC_URL)
	if err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("oidc.NewProvider")
	}
	a.tokenVerifier = oidcProvider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
	})
}

func (a *App) Initialize() {
	InitLog()
	log.Info().Str("source", "APP").Msg("initializing app")
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) Run(addr string, port string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Accept-Encoding", "Content-Range", "Content-Disposition", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})

	if port == "" {
		port = "8010"
	}

	listen := addr + ":" + port
	log.Info().Str("source", "APP").Msgf("app run %s", listen)

	if err := http.ListenAndServe(listen, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)); err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("http.ListenAndServe")
	}
}

func (a *App) InitializeRoutes() {
	a.Router.Use(a.LoggingMiddleware)
	a.Router.HandleFunc("/upload/{lang}/{ftype}", a.handleUpload).Methods("POST")
	a.Router.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir("/data"))))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
