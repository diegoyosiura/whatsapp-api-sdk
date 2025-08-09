package domain

// PhoneList is the envelope returned by GET /{WABA-ID}/phone_numbers.
type PhoneList struct {
	Data   []Phone `json:"data"`
	Paging *Paging `json:"paging,omitempty"`
}
