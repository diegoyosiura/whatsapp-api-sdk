package domain

type TemplateMessage struct {
	Template *TemplateBody `json:"template"`
}
type TemplateBody struct {
	Name       string               `json:"name"`
	Language   *TemplateLanguage    `json:"language"`
	Components []*TemplateComponent `json:"components"`
}
type TemplateLanguage struct {
	Code string `json:"code"`
}
type TemplateComponent struct {
	Type       string               `json:"type"`
	SubType    *string              `json:"sub_type"`
	Index      *string              `json:"index"`
	Parameters []*TemplateParameter `json:"parameters"`
}

type TemplateParameter struct {
	Type    string  `json:"type"`
	Text    *string `json:"text"`
	Payload *string `json:"payload"`
	*TemplateParameterCurrency
	*TemplateParameterDateTime
	*TemplateParameterImage
}

type TemplateParameterCurrency struct {
	Currency *TemplateParameterCurrencyOptions `json:"currency"`
}

type TemplateParameterCurrencyOptions struct {
	FallbackValue string `json:"fallback_value"`
	Code          string `json:"code"`
	Amount1000    int    `json:"amount_1000"`
}

type TemplateParameterDateTime struct {
	DateTime *TemplateParameterDateTimeOptions `json:"date_time"`
}

type TemplateParameterDateTimeOptions struct {
	FallbackValue string `json:"fallback_value"`
	DayOfWeek     int    `json:"day_of_week"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	DayOfMonth    int    `json:"day_of_month"`
	Hour          int    `json:"hour"`
	Minute        int    `json:"minute"`
	Calendar      string `json:"calendar"`
}

type TemplateParameterImage struct {
	Image *TemplateParameterImageOptions `json:"image"`
}

type TemplateParameterImageOptions struct {
	Link string `json:"link"`
}

func NewSendTemplateRequest(to, templateName, templateLang string, componentList []*TemplateComponent) *SendMessage {
	rt := "individual"

	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "template",
		TemplateMessage: &TemplateMessage{
			Template: &TemplateBody{
				Name: templateName,
				Language: &TemplateLanguage{
					Code: templateLang,
				},
				Components: componentList,
			},
		},
	}
}

func (s *SendMessage) validateTemplateMessage() error {
	return nil
}
