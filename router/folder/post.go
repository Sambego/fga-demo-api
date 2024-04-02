package folder

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

type CreateFolderResponse struct {
	Folder data.Folder `json:"folder"`
}

func parseCreateBody(body io.Reader) data.Folder {
	var parsedBody data.Folder
	err := json.NewDecoder(body).Decode(&parsedBody)

	if err != nil {
		log.Fatal("Somehthing went wrong parsing the request body.")
	}

	return parsedBody
}

func CreateFolderHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) (*CreateFolderResponse, error) {
	log.Printf("[POST] /documents ")

	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Create the new folder
	folder := store.CreateFolder(parseCreateBody(request.Body))

	// Create the ownership tupple for the current user
	tuples := []openfga.TupleKey{
		{
			Object:   fmt.Sprintf("folder:%s", folder.ID),
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
		return nil, err
	}

	// Return the created folder
	log.Printf("            -> Created folder with ID %v", folder.ID)
	return &CreateFolderResponse{
		Folder: folder,
	}, nil
}
