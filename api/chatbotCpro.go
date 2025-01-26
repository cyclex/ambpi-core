package api

type TextCpro struct {
	Body string `json:"body"`
}

type PayloadCpro struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Text             TextCpro `json:"text"`
}

type CproParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CproComponent struct {
	Type       string          `json:"type"`
	Parameters []CproParameter `json:"parameters"`
}

type CproLanguage struct {
	Code   string `json:"code"`
	Policy string `json:"policy"`
}

type CproTemplate struct {
	Name       string          `json:"name"`
	Language   CproLanguage    `json:"language"`
	Components []CproComponent `json:"components"`
}

type CproPayloadPush struct {
	MessagingProduct string       `json:"messaging_product"`
	To               string       `json:"to"`
	RecipientType    string       `json:"recipient_type"`
	Type             string       `json:"type"`
	Template         CproTemplate `json:"template"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	ExpiresAfter string `json:"expires_after"`
}
