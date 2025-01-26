package domain

import (
	"context"

	"github.com/cyclex/ambpi-core/api"
)

type CmsUcase interface {
	Login(c context.Context, req api.Login) (data map[string]interface{}, err error)
	CheckToken(c context.Context, req api.CheckToken) error
	Report(c context.Context, req api.Report, category string) (data map[string]interface{}, err error)
	FindDetailRedeem(c context.Context, category string) (data map[string]interface{}, err error)

	SendPushNotif(c context.Context, req api.SendPushNotif) (err error)
	GetProgram(c context.Context) (data []map[string]interface{}, err error)
	SetProgram(c context.Context, cond, data map[string]interface{}) (err error)

	ImportPrize(c context.Context, req api.Job) (status bool, totalRows int, err error)
	DownloadRedeem(c context.Context, req api.Job) (pathToFile string, status bool, totalRows int, err error)

	CreateJob(c context.Context, req api.Job) (err error)
	ListJob(c context.Context, category string) (data []map[string]interface{}, err error)

	CreateUser(c context.Context, req api.User) (err error)
	SetUser(c context.Context, req api.User) (err error)
	DeleteUser(c context.Context, deletedID int64) (err error)
	SetUserPassword(c context.Context, req api.User) (err error)

	ValidateRedeem(c context.Context, req api.ValidateRedeem) (err error)
}
