package api

// Define the structs to match the JSON payload
type CproMessageText struct {
	Body string `json:"body"`
}

type CproMessage struct {
	From      string           `json:"from"`
	ID        string           `json:"id"`
	Timestamp string           `json:"timestamp"`
	Text      CproMessageText  `json:"text"`
	Type      string           `json:"type"`
	Image     CproMessageImage `json:"image"`
}

type CproMessageImage struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	ID       string `json:"id"`
}

type CproContactProfile struct {
	Name string `json:"name"`
}

type CproContact struct {
	Profile CproContactProfile `json:"profile"`
	WaID    string             `json:"wa_id"`
}

type CproMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type CproValue struct {
	MessagingProduct string        `json:"messaging_product"`
	Metadata         CproMetadata  `json:"metadata"`
	Contacts         []CproContact `json:"contacts"`
	Messages         []CproMessage `json:"messages"`
}

type CproChange struct {
	Value CproValue `json:"value"`
	Field string    `json:"field"`
}

type CproEntry struct {
	ID      string       `json:"id"`
	Changes []CproChange `json:"changes"`
}

type CproWebhookPayload struct {
	Object string      `json:"object"`
	Entry  []CproEntry `json:"entry"`
}

type CproAccount struct {
	ID           string `json:"id"`
	Account      string `json:"account"`
	AccountTitle string `json:"account_title"`
}

type CproPayload struct {
	ID        string             `json:"id"`
	MID       string             `json:"mid"`
	ClientID  string             `json:"client_id"`
	ChannelID string             `json:"channel_id"`
	Account   CproAccount        `json:"account"`
	Data      CproWebhookPayload `json:"data"` // Using `any` since "data" is an empty object (can be changed to `map[string]interface{}` if needed)
}
