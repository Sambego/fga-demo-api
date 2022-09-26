package data

type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Parent string `json:"parent,omitempty"`
	Content string `json:"content"`
}

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Parent string `json:"parent,omitempty"`
}