package share

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	auth0fga "github.com/auth0-lab/fga-go-sdk"
	"github.com/gorilla/mux"
	"github.com/sambego/fga-demo-api/data"
)

type ShareDocumentRequestBody struct {
	Relation string `json:"relation"`
	User     string `json:"user"`
}

func parseShareBody(body io.Reader) ShareDocumentRequestBody {
	var parsedBody ShareDocumentRequestBody
	err := json.NewDecoder(body).Decode(&parsedBody)

	if err != nil {
		log.Fatal("Somehthing went wrong parsing the request body.")
	}

	return parsedBody
}

func ShareDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) error {

	// Get the shared document ID
	vars := mux.Vars(request)

	log.Printf("  - [POST] /documents/%v/share", vars["id"])

	// Parse the body
	newRelation := parseShareBody(request.Body)
	log.Printf("  - %v", newRelation.User)

	// Create the new relationship tupple for the current user
	tuples := []auth0fga.TupleKey{
		{
			Object:   auth0fga.PtrString(fmt.Sprintf("document:%s", vars["id"])),
			Relation: auth0fga.PtrString(newRelation.Relation),
			User:     auth0fga.PtrString(newRelation.User),
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
	log.Printf("            -> Shared the document with %v as a %v", newRelation.User, newRelation.Relation)

	// Error handling
	if err != nil {
		return err
	}

	return nil
}

func ShareFolderHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) error {
	// Get the shared document ID
	vars := mux.Vars(request)

	log.Printf("  - [POST] /folders/%v/share", vars["id"])

	// Parse the body
	newRelation := parseShareBody(request.Body)

	// Create the new relationship tupple for the current user
	tuples := []auth0fga.TupleKey{
		{
			Object:   auth0fga.PtrString(fmt.Sprintf("folder:%s", vars["id"])),
			Relation: auth0fga.PtrString(newRelation.Relation),
			User:     auth0fga.PtrString(newRelation.User),
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
	log.Printf("            -> Shared the folder with %v as a %v", newRelation.User, newRelation.Relation)

	// Error handling
	if err != nil {
		return err
	}

	return nil
}
