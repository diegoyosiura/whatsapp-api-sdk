package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

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
	if s.Type != "location" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be location", Op: "validateLocationMessage"}
	}

	if s.LocationMessage == nil {
		return &errorsx.ValidationError{Field: "LocationMessage", Reason: "location is nil", Op: "validateLocationMessage"}
	}

	if s.LocationMessage.Location == nil {
		return &errorsx.ValidationError{Field: "Location", Reason: "location is nil", Op: "validateLocationMessage"}
	}

	if s.LocationMessage.Location.Name == "" {
		return &errorsx.ValidationError{Field: "Location.Name", Reason: "empty string", Op: "validateLocationMessage"}
	}

	if s.LocationMessage.Location.Latitude == "" {
		return &errorsx.ValidationError{Field: "Location.Latitude", Reason: "empty string", Op: "validateLocationMessage"}
	}

	if s.LocationMessage.Location.Longitude == "" {
		return &errorsx.ValidationError{Field: "Location.Longitude", Reason: "empty string", Op: "validateLocationMessage"}
	}

	if s.LocationMessage.Location.Address == "" {
		return &errorsx.ValidationError{Field: "Location.Address", Reason: "empty string", Op: "validateLocationMessage"}
	}
	return nil
}
