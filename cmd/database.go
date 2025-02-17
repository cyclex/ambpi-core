package main

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(driver, dsn string, debug bool) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if debug {
		db = db.Debug()
	}

	db.AutoMigrate(
		&model.ConversationsLog{},
		&model.UsersUniqueCode{}, &model.Prizes{}, &model.RedeemPrizes{},
		&model.UserCMS{}, &model.Program{},
	)

	return db, nil
}

func ConnectQueue(dsn string, c context.Context) (*mongo.Client, error) {

	timeOutRequest := 30
	_, cancel := context.WithTimeout(c, time.Duration(timeOutRequest)*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	err = client.Connect(c)
	if err != nil {
		return nil, err
	}

	err = client.Ping(c, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
