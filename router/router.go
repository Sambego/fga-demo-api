package router

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sambego/fga-demo-api/middleware/auth"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func CreateRouter(service Service) *mux.Router {
	router := mux.NewRouter()

	tokenSecret, ok := os.LookupEnv("TOKEN_SECRET")
	if !ok {
		log.Fatal("'TOKEN_SECRET' environment variable must be set")
	}
	router.Use(auth.JWTTokenVerifierMiddleware(tokenSecret))

	router.HandleFunc("/", defaultHandler).Methods(http.MethodGet)
	router.HandleFunc("/documents", service.GetDocuments).Methods(http.MethodGet)
	router.HandleFunc("/documents", service.CreateDocument).Methods(http.MethodPost)
	router.HandleFunc("/documents/{id}", service.GetDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents/{id}/share", service.ShareDocument).Methods(http.MethodPost)
	router.HandleFunc("/folders", service.GetFolders).Methods(http.MethodGet)
	router.HandleFunc("/folders", service.CreateFolder).Methods(http.MethodPost)
	router.HandleFunc("/folders/{id}", service.GetFolder).Methods(http.MethodGet)
	router.HandleFunc("/folders/{id}/share", service.ShareFolder).Methods(http.MethodPost)

	return router
}
