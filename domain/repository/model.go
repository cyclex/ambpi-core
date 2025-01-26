package repository

import (
	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain/model"
)

type ModelRepository interface {
	CreateConversationsLog(new model.ConversationsLog) (err error)

	CreateUsersUniqueCodes(new model.UsersUniqueCode) (affected int64, err error, pk uint)

	FindFirstPrizes(cond map[string]interface{}) (data model.Prizes, err error)
	SetPrizes(cond, updated map[string]interface{}) (affected int64, err error)
	CreatePrize(rows [][]string) (totalRows int, err error)
	FindActivePrizes(cond map[string]interface{}, isActive bool) (data model.Prizes, err error)

	CreateRedeemPrizes(new model.RedeemPrizes) (affected int64, err error)
	FindFirstRedeemPrizes(cond map[string]interface{}) (data model.RedeemPrizes, err error)
	SetRedeemPrizes(cond, updated map[string]interface{}) (err error)
	FindRedeemID(id string) (data map[string]interface{}, err error)

	FindToken() (data model.Token, err error)
	SetToken(updated map[string]interface{}) (err error)

	Login(username, password string) (data model.UserCMS, err error)
	SetTokenLogin(id uint, token string) error
	CheckToken(token string) error

	ReportSummary() (data map[string]interface{}, err error)
	ReportHistoryValidation(req api.Report) (data map[string]interface{}, err error)
	ReportPrize(req api.Report) (data map[string]interface{}, err error)
	ReportRedeem(req api.Report) (data map[string]interface{}, err error)
	ReportUsage(req api.Report) (data map[string]interface{}, err error)
	ReportAvailability() (data map[string]interface{}, err error)

	FindProgram() (data []map[string]interface{}, err error)
	SetProgram(cond, updated map[string]interface{}) (affected int64, err error)
	IsProgramActive(retail string) (status int)

	FindUserBy(cond map[string]interface{}) (data []model.UserCMS, err error)
	SetUser(id int64, kol model.UserCMS) (err error)
	RemoveUser(id []int64) (err error)
	CreateUser(new model.UserCMS) (err error)
	SetUserPassword(username string, kol model.UserCMS) (err error)
}
