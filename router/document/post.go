package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	openfga "github.com/openfga/go-sdk"
	openfgaclient "github.com/openfga/go-sdk/client"
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

func CreateDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) (*CreateDocumentResponse, error) {
	log.Printf("[POST] /documents ")
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Create the new document
	document := store.CreateDocument(parseCreateBody(request.Body))

	// Create the ownership tupple for the current user
	tuples := []openfga.TupleKey{
		{
			Object:   fmt.Sprintf("doc:%s", document.ID),
			Relation: "owner",
			User:     fmt.Sprintf("user:%s", authCtx.Subject),
		},
	}

	// Create an FGA request
	body := openfgaclient.ClientWriteRequest{
		Writes: tuples,
	}

	// Makethe FGA request
	_, err := FGAClient.Write(request.Context()).Body(body).Execute()

	// Error handling
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Return the created document
	log.Printf("            -> Created document with ID %v", document.ID)
	return &CreateDocumentResponse{
		Document: document,
	}, nil
}
