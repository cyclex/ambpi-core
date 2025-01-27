package model

import "github.com/jinzhu/gorm"

type ConversationsLog struct {
	gorm.Model
	SessionID string `gorm:"session_id" json:"sessionID"`
	Incoming  string `gorm:"incoming" json:"waPayload"`
	Outgouing string `gorm:"outgoing" json:"chat"`
	WAID      string `gorm:"wa_id" json:"waID"`
}
