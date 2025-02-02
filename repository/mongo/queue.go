package mongo

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	DB            *mongo.Database
	c             context.Context
	channel       string
	expiredInHour time.Duration
}

func NewmongoRepository(c context.Context, db *mongo.Database, channel string, expired time.Duration) repository.QueueRepository {
	return &mongoRepo{
		DB:            db,
		c:             c,
		channel:       channel,
		expiredInHour: expired,
	}
}

func AckQueue(exp time.Duration) bson.D {

	data := bson.D{
		{"expired_at", time.Now().Local().Add(exp)},
		{"state", 2},
	}

	return data
}

func QueueBasedOnState(stateID int) bson.D {

	data := bson.D{
		{"state", stateID},
	}

	return data
}

func timeToLaunch(stateID int) bson.D {

	data := bson.D{
		{"state", stateID},
		{"messages.scheduledat", bson.D{{"$lt", time.Now().Local().Unix()}}},
	}

	return data
}

func (self *mongoRepo) CreateQueueRedeem(data model.QueueRedeem) error {

	data.ExpiredAt = time.Now().Local().Add(96 * time.Hour)

	_, err := self.DB.Collection("redeem").InsertOne(self.c, data)

	return err

}

func (self *mongoRepo) GetQueueRedeem() (queue []model.QueueRedeem, err error) {

	condition := QueueBasedOnState(1)

	dt, err := self.DB.Collection("redeem").Find(
		self.c,
		condition,
	)

	if err != nil {
		return nil, err
	}

	for dt.Next(self.c) {
		var queueHolder model.QueueRedeem

		err := dt.Decode(&queueHolder)
		if err != nil {
			return nil, err
		}

		queue = append(queue, queueHolder)
	}

	defer dt.Close(self.c)

	return queue, nil

}

func (self *mongoRepo) UpdateQueueRedeem(id primitive.ObjectID) error {

	id, err := primitive.ObjectIDFromHex(id.Hex())
	if err != nil {
		return err
	}

	col := self.DB.Collection("redeem")

	updatedData := AckQueue(self.expiredInHour * time.Hour)

	_, err = col.UpdateOne(
		self.c,
		bson.M{"_id": id},
		bson.D{
			{"$set", updatedData},
		})

	if err != nil {
		return err
	}

	return nil

}

func (self *mongoRepo) CreateQueueReply(data model.QueueReply) error {

	data.ExpiredAt = time.Now().Local().Add(96 * time.Hour)

	_, err := self.DB.Collection("reply").InsertOne(self.c, data)

	return err

}

func (self *mongoRepo) GetQueueReply() (queue []model.QueueReply, err error) {

	// condition := timeToLaunch(1)
	condition := QueueBasedOnState(1)

	dt, err := self.DB.Collection("reply").Find(
		self.c,
		condition,
	)

	if err != nil {
		return nil, err
	}

	for dt.Next(self.c) {
		var queueHolder model.QueueReply

		err := dt.Decode(&queueHolder)
		if err != nil {
			return nil, err
		}

		queue = append(queue, queueHolder)
	}

	defer dt.Close(self.c)

	return queue, nil

}

func (self *mongoRepo) UpdateQueueReply(id primitive.ObjectID) error {

	id, err := primitive.ObjectIDFromHex(id.Hex())
	if err != nil {
		return err
	}

	col := self.DB.Collection("reply")

	updatedData := AckQueue(self.expiredInHour * time.Hour)

	_, err = col.UpdateOne(
		self.c,
		bson.M{"_id": id},
		bson.D{
			{"$set", updatedData},
		})

	if err != nil {
		return err
	}

	return nil

}

func (self *mongoRepo) CreateJob(data model.QueueJob) error {

	data.ExpiredAt = time.Now().Local().Add(7 * (24 * time.Hour))
	data.CreatedAt = time.Now().Local().Format("2006-01-02 15:04:05")

	_, err := self.DB.Collection("jobs").InsertOne(self.c, data)

	return err

}

func (self *mongoRepo) GetJob(category string) ([]model.QueueJob, error) {

	// Define a slice to hold the results
	var queue []model.QueueJob

	// Perform the Find operation on the "jobs" collection
	cond := bson.M{}
	if category != "" {
		if category == "download" {
			cond = bson.M{"job_type": bson.M{"$in": []string{"download_redeem", "download_history_validation"}}}
		} else {
			cond = bson.M{"job_type": category}
		}
	}

	// Define the sort option to order by created_at descending
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := self.DB.Collection("jobs").Find(self.c, cond, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(self.c) // Ensure the cursor is closed when the function ends

	// Iterate through the cursor and decode each document into queueHolder
	for cursor.Next(self.c) {
		var queueHolder model.QueueJob
		if err := cursor.Decode(&queueHolder); err != nil {
			return nil, err
		}
		queue = append(queue, queueHolder)
	}

	// Check if any errors occurred during the iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return queue, nil
}

func (self *mongoRepo) UpdateJob(data model.QueueJob) (err error) {

	id, err := primitive.ObjectIDFromHex(data.ID.Hex())
	if err != nil {
		return err
	}

	col := self.DB.Collection("jobs")

	newExpiration := time.Now().Local().Add(24 * 7 * time.Hour)
	updatedData := bson.D{
		{"$set", bson.D{
			{"job_status", data.JobStatus},
			{"total_rows", data.TotalRows},
			{"file", data.File},
			{"expired_at", newExpiration},
		}},
	}

	_, err = col.UpdateOne(
		self.c,
		bson.M{"_id": id},
		updatedData)

	if err != nil {
		return err
	}

	return nil

}
