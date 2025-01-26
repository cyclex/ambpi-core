package api

// Request
type Text struct {
	Body string `json:"body"`
}

type Data struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}

type ReqSendMessageText struct {
	XID         string `json:"xid"`
	ChannelID   string `json:"channel_id"`
	AccountID   string `json:"account_id"`
	DivisionID  string `json:"division_id"`
	IsHelpdesk  bool   `json:"is_helpdesk"`
	MessageType string `json:"message_type"`
	Data        Data   `json:"data"`
}
