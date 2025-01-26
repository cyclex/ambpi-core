package model

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	AccessToken string `gorm:"access_token" json:"accessToken"`
	ExpiredAt   string `gorm:"expired_at" json:"expiredAt"`
}
