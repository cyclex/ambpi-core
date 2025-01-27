package usecase

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ordersUcase struct {
	q              repository.QueueRepository
	contextTimeout time.Duration
}

func NewOrdersUcase(q repository.QueueRepository, timeout time.Duration) domain.OrdersUcase {

	return &ordersUcase{
		q:              q,
		contextTimeout: timeout,
	}
}

func (self *ordersUcase) CreateQueueRedeem(c context.Context, msg api.PayloadRedeem) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data := model.QueueRedeem{
		TrxId:    uuid.New().String(),
		State:    1,
		Messages: msg,
	}

	err = self.q.CreateQueueRedeem(data)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateQueueRedeem] CreateQueueRedeem")
	}

	return

}

func (self *ordersUcase) GetQueueRedeem(c context.Context) (data []model.QueueRedeem, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err = self.q.GetQueueRedeem()
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetQueueRedeem] GetQueueRedeem")
	}

	return

}

func (self *ordersUcase) UpdateQueueRedeem(c context.Context, id primitive.ObjectID) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.q.UpdateQueueRedeem(id)
	if err != nil {
		err = errors.Wrap(err, "[usecase.UpdateQueueRedeem] UpdateQueueRedeem")
	}

	return
}

func (self *ordersUcase) GetQueueReply(c context.Context) (data []model.QueueReply, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err = self.q.GetQueueReply()
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetQueueReply] GetQueueReply")
	}

	return

}

func (self *ordersUcase) CreateQueueReply(c context.Context, msg api.PayloadReply) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data := model.QueueReply{
		TrxId:    uuid.New().String(),
		State:    1,
		Messages: msg,
	}

	err = self.q.CreateQueueReply(data)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateQueueReply] CreateQueueReply")
	}

	return

}

func (self *ordersUcase) UpdateQueueReply(c context.Context, id primitive.ObjectID) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.q.UpdateQueueReply(id)
	if err != nil {
		err = errors.Wrap(err, "[usecase.UpdateQueueReply] UpdateQueueReply")
	}

	return

}

func (self *ordersUcase) GetJob(c context.Context, category string) (data []model.QueueJob, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err = self.q.GetJob(category)
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetJob] GetJob")
	}

	return
}

func (self *ordersUcase) UpdateJob(c context.Context, data api.Job) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.q.UpdateJob(model.QueueJob{ID: data.ID, JobStatus: data.JobStatus, TotalRows: data.TotalRows, File: data.File})
	if err != nil {
		err = errors.Wrap(err, "[usecase.UpdateJob] UpdateJob")
	}

	return

}
