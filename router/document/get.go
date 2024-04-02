package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	openfgaclient "github.com/openfga/go-sdk/client"
	"github.com/sambego/fga-demo-api/data"
	"github.com/sambego/fga-demo-api/errors"
	"github.com/sambego/fga-demo-api/middleware/auth"
)

type GetDocumentResponse struct {
	Document data.Document `json:"document"`
}

func GetDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) (*GetDocumentResponse, error) {
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Get the document ID from the request variables
	vars := mux.Vars(request)

	log.Printf("[GET] /documents/%v", vars["id"])

	// Create an FGA request to check of the current user can view the document
	body := openfgaclient.ClientCheckRequest{
		Object:   fmt.Sprintf("doc:%s", vars["id"]),
		Relation: "can_read",
		User:     fmt.Sprintf("user:%s", authCtx.Subject),
	}

	// Makethe FGA request
	resp, err := FGAClient.Check(request.Context()).Body(body).Execute()

	// Error Handling
	if err != nil {
		// handle error
	}

	// Return a forbidden response when the current user is not allowed to view the document
	if !resp.GetAllowed() {
		log.Printf("            -> Not allowed for user %v", authCtx.Subject)
		return nil, errors.ErrorForbidden
	}

	// Return the document if we're allowed to
	document := store.GetDocument(vars["id"])
	log.Printf("            -> Returned document with name %v", document.Name)
	return &GetDocumentResponse{
		Document: document,
	}, nil
}

// For debugging purposes, no FGA checks are done here
func GetDocumentsHandler(writer http.ResponseWriter, request *http.Request, store *data.Store) {
	log.Printf("[GET] /documents ")
	document, _ := json.Marshal(store.GetDocuments())
	io.WriteString(writer, string(document))
	log.Printf("            -> Returned all documents")
}
