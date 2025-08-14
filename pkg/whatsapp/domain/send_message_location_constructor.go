package domain

type LocationMessage struct {
	Location *LocationBody `json:"location"`
}
type LocationBody struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Address   string `json:"address"`
}

func NewSendLocationRequest(to, latitude, longitude, name, address string) *SendMessage {
	rt := "individual"

	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "location",
		LocationMessage: &LocationMessage{
			Location: &LocationBody{
				Latitude:  latitude,
				Longitude: longitude,
				Name:      name,
				Address:   address,
			},
		},
	}
}

func NewSendContextLocationRequest(to, latitude, longitude, name, address, targetMessage string) *SendMessage {
	rt := "individual"

	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "location",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		LocationMessage: &LocationMessage{
			Location: &LocationBody{
				Latitude:  latitude,
				Longitude: longitude,
				Name:      name,
				Address:   address,
			},
		},
	}
}
func (s *SendMessage) validateLocationMessage() error {
	return nil
}
