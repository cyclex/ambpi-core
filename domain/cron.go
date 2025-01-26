package domain

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersUcase interface {
	GetQueueRedeem(c context.Context) (data []model.QueueRedeem, err error)
	UpdateQueueRedeem(c context.Context, id primitive.ObjectID) (err error)
	CreateQueueRedeem(c context.Context, data api.PayloadRedeem) (err error)

	CreateQueueReply(c context.Context, data api.PayloadReply) (err error)
	GetQueueReply(c context.Context) (data []model.QueueReply, err error)
	UpdateQueueReply(c context.Context, id primitive.ObjectID) (err error)

	GetJob(c context.Context, category string) (data []model.QueueJob, err error)
	UpdateJob(c context.Context, data api.Job) (err error)
}

type CronChatbot struct {
	ID              string
	Err             error
	TrxChatBotMsgID string
	ServerTime      time.Time
	Status          string
}
