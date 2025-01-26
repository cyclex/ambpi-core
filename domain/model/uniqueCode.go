package model

import "github.com/jinzhu/gorm"

type UsersUniqueCode struct {
	gorm.Model
	WaID          string `gorm:"wa_id" json:"waID"`
	Program       string `gorm:"program" json:"program"`
	Name          string `gorm:"name" json:"name"`
	Profession    string `gorm:"profession" json:"profession"`
	SessionID     string `gorm:"session_id" json:"session_id"`
	IsFormatValid bool   `gorm:"is_format_valid" json:"isFormatValid"`
	Msisdn        string `gorm:"msisdn" json:"msisdn"`
	Raw           string `gorm:"raw" json:"raw"`
	NIK           string `gorm:"nik" json:"nik"`
	IsZonk        bool   `gorm:"is_zonk" json:"isZonk"`
	Reply         string `gorm:"reply" json:"reply"`
	County        string `gorm:"county" json:"county"`
}

type Prizes struct {
	gorm.Model
	Prize          string `gorm:"prize" json:"prize"`
	IsUsed         bool   `gorm:"is_used" json:"isUsed"`
	PrizeType      string `gorm:"prize_type" json:"type"`
	SequenceNumber int    `gorm:"sequence_number" json:"secuenceNumber"`
}

type RedeemPrizes struct {
	gorm.Model
	PrizeID           uint   `gorm:"prize_id" json:"prizeID"`
	UsersUniqueCodeID uint   `gorm:"users_unique_code_id" json:"usersUniqueCodeID"`
	PushID            uint   `gorm:"push_id" json:"pushID"`
	Msisdn            string `gorm:"msisdn" json:"msisdn"`
	Notes             string `gorm:"notes" json:"notes"`
	Approved          bool   `gorm:"approved" json:"approved"`
	Amount            int    `gorm:"amount" json:"amount"`
	Author            string `gorm:"author" json:"author"`
	LotteryNumber     string `gorm:"lottery_number" json:"lotteryNumber"`
	DateValidation    string `gorm:"date_validation" json:"dateValidation"`
}
