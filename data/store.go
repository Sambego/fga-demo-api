package data

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Store struct {
	Documents []Document
	Folders   []Folder
}

/*
 * Documents
 */
func (store *Store) CreateDocument(document Document) Document {
	id := uuid.NewString()
	document.ID = id
	store.Documents = append(store.Documents, document)

	return document
}

func (store *Store) GetDocument(id string) Document {
	idx := slices.IndexFunc(store.Documents, func(document Document) bool { return document.ID == id })

	return store.Documents[idx]
}

func (store *Store) GetDocuments() []Document {
	return store.Documents
}

/*
 * Folder
 */
func (store *Store) CreateFolder(folder Folder) Folder {
	id := uuid.NewString()
	folder.ID = id
	store.Folders = append(store.Folders, folder)

	return folder
}

func (store *Store) GetFolder(id string) Folder {
	idx := slices.IndexFunc(store.Folders, func(folder Folder) bool { return folder.ID == id })

	return store.Folders[idx]
}

func (store *Store) GetFodlers() []Folder {
	return store.Folders
}
