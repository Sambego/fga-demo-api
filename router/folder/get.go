package folder

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

type GetFolderResponse struct {
	Folder data.Folder `json:"folder"`
}

func GetFolderHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) (*GetFolderResponse, error) {
	// Get the Authentication context
	authCtx, _ := auth.AuthContextFromContext(request.Context())

	// Get the document ID from the request variables
	vars := mux.Vars(request)

	log.Printf("  - [GET] /folders/%v", vars["id"])

	// Create an FGA request to check of the current user can view the folder
	body := auth0fga.CheckRequest{
		TupleKey: &auth0fga.TupleKey{
			Object:   auth0fga.PtrString(fmt.Sprintf("folder:%s", vars["id"])),
			Relation: auth0fga.PtrString("viewer"),
			User:     auth0fga.PtrString(authCtx.Subject),
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

	// Return the folder if we're allowed to
	folder := store.GetFolder(vars["id"])
	log.Printf("            -> Returned folder with name %v", folder.Name)
	return &GetFolderResponse{
		Folder: folder,
	}, nil
}

// For debuggin purposes, no FGA checks are done here
func GetFoldersHandler(writer http.ResponseWriter, request *http.Request, store *data.Store, FGAClient auth0fga.Auth0FgaApi) {
	log.Printf("  - [GET] /folders ")
	folder, _ := json.Marshal(store.GetFodlers())
	io.WriteString(writer, string(folder))
	log.Printf("            -> Returned all folders")
}
