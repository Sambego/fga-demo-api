package share

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	openfga "github.com/openfga/go-sdk"
	openfgaclient "github.com/openfga/go-sdk/client"
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

func ShareDocumentHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) error {

	// Get the shared document ID
	vars := mux.Vars(request)

	log.Printf("[POST] /documents/%v/share", vars["id"])

	// Parse the body
	newRelation := parseShareBody(request.Body)

	// Create the new relationship tupple for the current user
	tuples := []openfga.TupleKey{
		{
			Object:   fmt.Sprintf("doc:%s", vars["id"]),
			Relation: newRelation.Relation,
			User:     newRelation.User,
		},
	}

	// Create an FGA request
	body := openfgaclient.ClientWriteRequest{
		Writes: tuples,
	}

	// Makethe FGA request
	_, err := FGAClient.Write(request.Context()).Body(body).Execute()
	log.Printf("            -> Shared the document with %v as a %v", newRelation.User, newRelation.Relation)

	// Error handling
	if err != nil {
		return err
	}

	return nil
}

func ShareFolderHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) error {
	// Get the shared document ID
	vars := mux.Vars(request)

	log.Printf("[POST] /folders/%v/share", vars["id"])

	// Parse the body
	newRelation := parseShareBody(request.Body)

	// Create the new relationship tupple for the current user
	tuples := []openfga.TupleKey{
		{
			Object:   fmt.Sprintf("folder:%s", vars["id"]),
			Relation: newRelation.Relation,
			User:     newRelation.User,
		},
	}

	// Create an FGA request
	body := openfgaclient.ClientWriteRequest{
		Writes: tuples,
	}

	// Makethe FGA request
	_, err := FGAClient.Write(request.Context()).Body(body).Execute()
	log.Printf("            -> Shared the folder with %v as a %v", newRelation.User, newRelation.Relation)

	// Error handling
	if err != nil {
		return err
	}

	return nil
}
