package router

import (
	"encoding/json"
	"net/http"

	auth0fga "github.com/auth0-lab/fga-go-sdk"
	"github.com/sambego/fga-demo-api/data"
	"github.com/sambego/fga-demo-api/errors"
	"github.com/sambego/fga-demo-api/router/document"
	"github.com/sambego/fga-demo-api/router/folder"
	"github.com/sambego/fga-demo-api/router/share"
)

type Service struct {
	Store     *data.Store
	FGAClient auth0fga.Auth0FgaApi
}

/*
 * Documents
 */
func (service Service) CreateDocument(writer http.ResponseWriter, request *http.Request) {
	document, err := document.CreateDocumentHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Marshal document to JSON
	body, err := json.Marshal(document)
	if err != nil {
		http.Error(writer, "Failed to json marshal response", http.StatusInternalServerError)
		return
	}

	// Return the created document as JSON
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, "Oops, something went wrong!", http.StatusInternalServerError)
	}
}

func (service Service) GetDocument(writer http.ResponseWriter, request *http.Request) {
	document, err := document.GetDocumentHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Marshal document to JSON
	body, err := json.Marshal(document)
	if err != nil {
		http.Error(writer, "Failed to json marshal response", http.StatusInternalServerError)
		return
	}

	// Return the fetched document as JSON
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, "Oops, something went wrong!", http.StatusInternalServerError)
	}
}

func (service Service) ShareDocument(writer http.ResponseWriter, request *http.Request) {
	err := share.ShareDocumentHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

// For debuggin purposes, no FGA checks are done here
func (service Service) GetDocuments(writer http.ResponseWriter, request *http.Request) {
	document.GetDocumentsHandler(writer, request, service.Store, service.FGAClient)
}

/*
 * Folders
 */
func (service Service) CreateFolder(writer http.ResponseWriter, request *http.Request) {
	folder, err := folder.CreateFolderHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Marshal document to JSON
	body, err := json.Marshal(folder)
	if err != nil {
		http.Error(writer, "Failed to json marshal response", http.StatusInternalServerError)
		return
	}

	// Return the created document as JSON
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, "Oops, something went wrong!", http.StatusInternalServerError)
	}
}

func (service Service) GetFolder(writer http.ResponseWriter, request *http.Request) {
	folder, err := folder.GetFolderHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Marshal document to JSON
	body, err := json.Marshal(folder)
	if err != nil {
		http.Error(writer, "Failed to json marshal response", http.StatusInternalServerError)
		return
	}

	// Return the fetched document as JSON
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, "Oops, something went wrong!", http.StatusInternalServerError)
	}
}

func (service Service) ShareFolder(writer http.ResponseWriter, request *http.Request) {
	err := share.ShareFolderHandler(writer, request, service.Store, service.FGAClient)

	// Error handling
	if err != nil {
		if err == errors.ErrorUnauthorized {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

// For debuggin purposes, no FGA checks are done here
func (service Service) GetFolders(writer http.ResponseWriter, request *http.Request) {
	folder.GetFoldersHandler(writer, request, service.Store, service.FGAClient)
}
