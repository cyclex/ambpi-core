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

// func main() {
// 	http.HandleFunc("/webhook", webhookHandler)

// 	fmt.Println("Starting server on :8080...")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
