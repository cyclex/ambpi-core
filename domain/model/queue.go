package model

import (
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueRedeem struct {
	gorm.Model
	ID        primitive.ObjectID `bson:"_id"`
	TrxId     string             `bson:"trx_id"`
	Messages  api.PayloadRedeem  `bson:"messages"`
	State     int                `bson:"state"`
	ExpiredAt time.Time          `bson:"expired_at"`
}

type QueueReply struct {
	gorm.Model
	ID        primitive.ObjectID `bson:"_id"`
	TrxId     string             `bson:"trx_id"`
	Messages  api.PayloadReply   `bson:"messages"`
	State     int                `bson:"state"`
	ExpiredAt time.Time          `bson:"expired_at"`
	Raw       string             `bson:"raw"`
	IsPush    bool               `bson:"is_push"`
}

type QueueJob struct {
	gorm.Model
	ID        primitive.ObjectID `bson:"_id"`
	JobType   string             `bson:"job_type"`
	JobStatus string             `bson:"job_status"`
	Author    string             `bson:"author"`
	File      string             `bson:"file"`
	TotalRows int                `bson:"total_rows"`
	ExpiredAt time.Time          `bson:"expired_at"`
	StartAt   string             `bson:"start_at"`
	EndAt     string             `bson:"end_at"`
	CreatedAt string             `bson:"created_at"`
}
