package api

type ResponseChatbot struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	ServerTime int64       `json:"servertime"`
	Data       interface{} `json:"data"`
}
