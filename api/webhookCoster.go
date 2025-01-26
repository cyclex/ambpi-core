package api

type Profile struct {
	Name string `json:"name"`
}

type Contact struct {
	WAID    string  `json:"wa_id"`
	Profile Profile `json:"profile"`
}

type Message struct {
	From string `json:"from"`
	ID   string `json:"id"`
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type Value struct {
	Contacts         []Contact `json:"contacts"`
	MessagingProduct string    `json:"messaging_product"`
	Messages         []Message `json:"messages"`
	Metadata         struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberID      string `json:"phone_number_id"`
	} `json:"metadata"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type DataWebhook struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Account struct {
	ID           string `json:"id"`
	Account      string `json:"account"`
	AccountTitle string `json:"account_title"`
}

type PayloadWebhook struct {
	ID        string      `json:"id"`
	MID       string      `json:"mid"`
	ClientID  string      `json:"client_id"`
	ChannelID string      `json:"channel_id"`
	Account   Account     `json:"account"`
	Data      DataWebhook `json:"data"`
}
