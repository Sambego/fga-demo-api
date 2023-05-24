package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	auth0fga "github.com/auth0-lab/fga-go-sdk"
	"github.com/sambego/fga-demo-api/data"
	"github.com/sambego/fga-demo-api/middleware/auth"
)

type CreateDocumentResponse struct {
	Document data.Document `json:"document"`
}

func parseCreateBody(body io.Reader) data.Document {
	var parsedBody data.Document
	err := json.NewDecoder(body).Decode(&parsedBody)

	if err != nil {
		log.Fatal("Somehthing went wrong parsing the request body.")
	}

	return parsedBody
}

func CreateDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) (*CreateDocumentResponse, error) {
	log.Printf("  - [POST] /documents ")
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Create the new document
	document := store.CreateDocument(parseCreateBody(request.Body))

	// Create the ownership tupple for the current user
	tuples := []auth0fga.TupleKey{
		{
			Object:   auth0fga.PtrString(fmt.Sprintf("document:%s", document.ID)),
			Relation: auth0fga.PtrString("owner"),
			User:     auth0fga.PtrString(fmt.Sprintf("user:%s", authCtx.Subject)),
		},
	}

	// Create an FGA request
	body := auth0fga.WriteRequest{
		Writes: &auth0fga.TupleKeys{
			TupleKeys: tuples,
		},
	}

	// Makethe FGA request
	_, _, err := FGAClient.Write(request.Context()).Body(body).Execute()

	// Error handling
	if err != nil {
		return nil, err
	}

	// Return the created document
	log.Printf("            -> Created document with ID %v", document.ID)
	return &CreateDocumentResponse{
		Document: document,
	}, nil
}
