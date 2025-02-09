package api

type ResSendMessage struct {
	Messages []Message `json:"messages"`
	Err      []ErrDesc `json:"errors"`
}

type ErrDesc struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Image struct {
	Link    string `json:"link" bson:"link"`
	Caption string `json:"caption" bson:"caption"`
}

type ReqSendMessageImage struct {
	RecipientType string `json:"recipient_type"`
	To            string `json:"to"`
	Type          string `json:"type"`
	Image         Image  `json:"image"`
}

type ReqSendBroadcast struct {
	To       string   `json:"to"`
	Type     string   `json:"type"`
	Template Template `json:"template"`
	// Hsm      Hsm      `json:"hsm"`
}

type Template struct {
	Namespace  string       `json:"namespace"`
	Name       string       `json:"name"`
	Language   Language     `json:"language"`
	Components []Components `json:"components"`
}

type Hsm struct {
	Namespace   string       `json:"namespace"`
	ElementName string       `json:"element_name"`
	Language    Language     `json:"language"`
	LocalParam  []LocalParam `json:"localizable_params"`
}

type LocalParam struct {
	Default string `json:"default"`
}

type Language struct {
	Code   string `json:"code"`
	Policy string `json:"policy"`
}

type Components struct {
	Type       string       `json:"type"`
	Parameters []Parameters `json:"parameters"`
}

type Parameters struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResLogin struct {
	Users []User `json:"users"`
}

type PayloadRedeem struct {
	WaID          string `json:"waID"`
	Name          string `json:"name"`
	Profession    string `json:"profession"`
	IsFormatValid bool   `json:"isFormatValid"`
	Raw           string `json:"raw"`
	NIK           string `json:"nik"`
	County        string `json:"county"`
	MediaID       string `json:"mediaID"`
}

type PayloadReply struct {
	WaID        string   `json:"waID"`
	Chat        string   `json:"chat"`
	ScheduledAt int64    `json:"scheduledAt"`
	Raw         string   `json:"raw"`
	IsPush      bool     `json:"is_push"`
	Param       []string `json:"param"`
}
