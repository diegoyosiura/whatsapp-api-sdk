package domain

type PagingCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}
