package domain

// Paging mirrors Graph API paging metadata.
type Paging struct {
	Cursors  *PagingCursors `json:"cursors,omitempty"`
	Next     string         `json:"next,omitempty"`
	Previous string         `json:"previous,omitempty"`
}
