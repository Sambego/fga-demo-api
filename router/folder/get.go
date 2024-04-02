package folder

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

type GetFolderResponse struct {
	Folder data.Folder `json:"folder"`
}

func GetFolderHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient *openfgaclient.OpenFgaClient) (*GetFolderResponse, error) {
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Get the document ID from the request variables
	vars := mux.Vars(request)

	log.Printf("[GET] /folders/%v", vars["id"])

	// Create an FGA request to check of the current user can view the folder
	body := openfgaclient.ClientCheckRequest{
		Object:   fmt.Sprintf("folder:%s", vars["id"]),
		Relation: "viewer",
		User:     fmt.Sprintf("user:%s", authCtx.Subject),
	}

	// Makethe FGA request
	resp, err := FGAClient.Check(request.Context()).Body(body).Execute()

	// Error Handling
	if err != nil {
		// handle error
	}

	// Return forbidden response when the current user is not allowed to view the document
	if !resp.GetAllowed() {
		log.Printf("            -> Not allowed for user %v", authCtx.Subject)
		return nil, errors.ErrorForbidden
	}

	// Return the folder if we're allowed to
	folder := store.GetFolder(vars["id"])
	log.Printf("            -> Returned folder with name %v", folder.Name)
	return &GetFolderResponse{
		Folder: folder,
	}, nil
}

// For debugging purposes, no FGA checks are done here
func GetFoldersHandler(writer http.ResponseWriter, request *http.Request, store *data.Store) {
	log.Printf("[GET] /folders ")
	folder, _ := json.Marshal(store.GetFodlers())
	io.WriteString(writer, string(folder))
	log.Printf("            -> Returned all folders")
}
