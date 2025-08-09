package domain

type MessageContext struct {
	From string `json:"from,omitempty"`
	ID   string `json:"id,omitempty"`
}
