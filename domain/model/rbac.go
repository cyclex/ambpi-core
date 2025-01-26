package model

import "github.com/jinzhu/gorm"

type AccessCMS struct {
	gorm.Model
	PrivilegeID  uint `gorm:"privilege_id" json:"privilegeID"`
	RuleModulsID uint `gorm:"rule_moduls_id" json:"ruleModulsID"`
	HasAccess    bool `gorm:"has_access" json:"hasAccess"`
}

type Modul struct {
	gorm.Model
	Name       string `gorm:"name" json:"title"`
	Flag       bool   `gorm:"flag" json:"flag"`
	Controller string `gorm:"controller" json:"path"`
	OrderID    int    `gorm:"order_id" json:"orderID"`
	PID        int    `gorm:"pid" json:"pid"`
}

type Privilege struct {
	gorm.Model
	Name string `gorm:"name" json:"name"`
	Flag bool   `gorm:"flag" json:"flag"`
}

type RuleModul struct {
	gorm.Model
	ModulID uint `gorm:"modul_id" json:"modulID"`
	RuleID  uint `gorm:"rule_id" json:"ruleID"`
	Flag    bool `gorm:"flag" json:"flag"`
}

type Rule struct {
	gorm.Model
	Rule string `gorm:"rule" json:"rule"`
	Flag bool   `gorm:"flag" json:"flag"`
}

type UserCMS struct {
	gorm.Model
	Username  string `gorm:"username" json:"username"`
	Password  string `gorm:"password" json:"-"`
	Flag      bool   `gorm:"flag" json:"-"`
	Level     string `gorm:"level" json:"level"`
	Token     string `gorm:"token" json:"token"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}
