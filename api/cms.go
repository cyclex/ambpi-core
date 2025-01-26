package api

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"pass" validate:"required"`
}

type CheckToken struct {
	Token string `json:"token" validate:"required"`
}

type Access struct {
	UserID      string `json:"user_id"`
	PrivilegeID string `json:"privilege_id" validate:"required"`
}

type Report struct {
	From      string `json:"from" validate:"required"`
	To        string `json:"to" validate:"required"`
	PrizeType string `json:"prize_type"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Column    string `json:"column"`
	Keyword   string `json:"keyword"`
	Sort      string `json:"sort"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
	Level    string `json:"level"`
}

type ValidateRedeem struct {
	ID       int64  `json:"id" validate:"required"`
	Amount   int64  `json:"amount" validate:"required"`
	Notes    string `json:"notes"`
	Approved bool   `json:"approved"`
	Author   string `json:"author" validate:"required"`
}

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type SendPushNotif struct {
	RedeemID     uint     `json:"redeem_id" validate:"required"`
	PushBy       string   `json:"push_by" validate:"required"`
	TemplateName string   `json:"template_name" validate:"required"`
	Param        []string `json:"param"`
}

type Import struct {
	File string `json:"file" validate:"required"`
}

type ResponseSuccess struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type ResponseReport struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type SetProgram struct {
	Status    bool `json:"status"`
	StartDate int  `json:"startDate"`
	EndDate   int  `json:"endDate"`
}

type Prize struct {
	Prize       string `json:"prize" validate:"required"`
	Quota       int    `json:"quota" validate:"required"`
	ActiveQuota int    `json:"active_quota" validate:"required"`
	PrizeOwner  string `json:"prize_owner" validate:"required"`
	OpenAt      int    `json:"open_at" validate:"required"`
}

type Job struct {
	ID        primitive.ObjectID `json:"id"`
	File      string             `json:"file"`
	JobType   string             `json:"job_type" validate:"required"`
	Author    string             `json:"author" validate:"required"`
	TotalRows int                `json:"total_rows"`
	StartDate string             `json:"start_date"`
	EndDate   string             `json:"end_date"`
	JobStatus string             `json:"job_status"`
}
