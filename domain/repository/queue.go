package repository

import (
	"github.com/cyclex/ambpi-core/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueRepository interface {
	GetQueueRedeem() ([]model.QueueRedeem, error)
	UpdateQueueRedeem(id primitive.ObjectID) error
	CreateQueueRedeem(data model.QueueRedeem) error

	CreateQueueReply(data model.QueueReply) (err error)
	GetQueueReply() (data []model.QueueReply, err error)
	UpdateQueueReply(id primitive.ObjectID) (err error)

	CreateJob(data model.QueueJob) (err error)
	GetJob(category string) ([]model.QueueJob, error)
	UpdateJob(data model.QueueJob) (err error)
}
