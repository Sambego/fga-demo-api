package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	auth0fga "github.com/auth0-lab/fga-go-sdk"
	"github.com/gorilla/mux"
	"github.com/sambego/fga-demo-api/data"
	"github.com/sambego/fga-demo-api/errors"
	"github.com/sambego/fga-demo-api/middleware/auth"
)

type GetDocumentResponse struct {
	Document data.Document `json:"document"`
}

func GetDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) (*GetDocumentResponse, error) {
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Get the document ID from the request variables
	vars := mux.Vars(request)

	log.Printf("  - [GET] /documents/%v", vars["id"])

	// Create an FGA request to check of the current user can view the document
	body := auth0fga.CheckRequest{
		TupleKey: auth0fga.TupleKey{
			Object:   auth0fga.PtrString(fmt.Sprintf("document:%s", vars["id"])),
			Relation: auth0fga.PtrString("viewer"),
			User:     auth0fga.PtrString(fmt.Sprintf("user:%s", authCtx.Subject)),
		},
	}

	// Makethe FGA request
	resp, _, err := FGAClient.Check(request.Context()).Body(body).Execute()

	// Error Handling
	if err != nil {
		// handle error
	}

	// Return an unauthorized response when the current user is not allowed to view the document
	if !resp.GetAllowed() {
		log.Printf("            -> Not allowed for user %v", authCtx.Subject)
		return nil, errors.ErrorUnauthorized
	}

	// Return the document if we're allowed to
	document := store.GetDocument(vars["id"])
	log.Printf("            -> Returned document with name %v", document.Name)
	return &GetDocumentResponse{
		Document: document,
	}, nil
}

// For debuggin purposes, no FGA checks are done here
func GetDocumentsHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) {
	log.Printf("  - [GET] /documents ")
	document, _ := json.Marshal(store.GetDocuments())
	io.WriteString(writer, string(document))
	log.Printf("            -> Returned all documents")
}
