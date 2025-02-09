package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	_HttpDelivery "github.com/cyclex/ambpi-core/delivery/http"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/cyclex/ambpi-core/repository/mongo"
	"github.com/cyclex/ambpi-core/repository/postgre"
	"github.com/cyclex/ambpi-core/usecase"
)

func run_server(server, config string, debug bool) (err error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[run_server] panic occurred: %v", err)
		}
	}()

	// Create a context that can be cancelled when a shutdown signal is received
	c, cancel := context.WithCancel(context.Background())

	// Handle SIGINT (Ctrl+C) to initiate a graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v\n", sig)
		cancel() // Cancel the context to initiate graceful shutdown
	}()

	// load config
	cfg, err := pkg.LoadServiceConfig(config)
	if err != nil {

		err = errors.Wrap(err, "[run_server]")
		return
	}

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPass := cfg.Database.Pass
	dbName := cfg.Database.Name
	dbSsl := cfg.Database.Ssl
	dbTimeout := cfg.Database.Timeout

	if dbSsl == "" {
		dbSsl = "disable"
	}

	if dbTimeout <= 0 {
		dbTimeout = 5
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d", dbHost, dbPort, dbUser, dbName, dbPass, dbSsl, dbTimeout)
	conn, err := ConnectDB("postgres", dsn, debug)
	if err != nil {
		log.Fatal(err)
	}

	queueHost := cfg.Queue.Host
	queuePort := cfg.Queue.Port
	queueName := cfg.Queue.Name
	dsn = fmt.Sprintf("mongodb://%s:%d", queueHost, queuePort)
	queue, err := ConnectQueue(dsn, c)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = queue.Disconnect(c); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	timeoutCtx := time.Duration(30) * time.Second

	expired := cfg.Queue.Expired
	if expired < 1 {
		expired = 24
	}
	expiredQueue := time.Duration(expired)
	ordersQueue := mongo.NewmongoRepository(c, queue.Database(queueName), "chatbot", expiredQueue)

	phoneID := cfg.Chatbot.PhoneID
	urlSendMsg := cfg.Chatbot.Host
	urlPush := cfg.Chatbot.HostPush
	accessToken := cfg.Chatbot.AccessToken
	downloadFolder := cfg.DownloadFolder
	accountID := cfg.Chatbot.AccountID
	divisionID := cfg.Chatbot.DivisionID
	accessTokenPush := cfg.Chatbot.AccessTokenPush
	model := postgre.NewPostgreRepository(c, conn)
	chatUcase := usecase.NewChatUcase(model, urlPush, urlSendMsg, phoneID, accessToken, accountID, divisionID, accessTokenPush, downloadFolder, ordersQueue, rdb)
	cmsUcase := usecase.NewCmsUcase(model, timeoutCtx, chatUcase, ordersQueue, downloadFolder)
	orderUcase := usecase.NewOrdersUcase(ordersQueue, timeoutCtx)

	InitCron(orderUcase, chatUcase, cmsUcase, timeoutCtx, debug)
	e := echo.New()
	_HttpDelivery.NewCmsHandler(e, cmsUcase, debug)

	go func() {
		if err := e.Start(server); err != nil {
			log.Fatalf("[run_server] %s", err)
		}
	}()

	// Wait for the context to be cancelled (e.g., by receiving SIGINT)
	<-c.Done()

	log.Println("Shutting down gracefully...")

	// Create a context with a timeout for shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Shutdown the server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped.")

	return nil
}

func run_webhook(server, config string, debug bool) (err error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[run_webhook] panic occurred: %v", err)
		}
	}()

	// Create a context that can be cancelled when a shutdown signal is received
	c, cancel := context.WithCancel(context.Background())

	// Handle SIGINT (Ctrl+C) to initiate a graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v\n", sig)
		cancel() // Cancel the context to initiate graceful shutdown
	}()

	// load config
	cfg, err := pkg.LoadServiceConfig(config)
	if err != nil {
		err = errors.Wrap(err, "[run_webhook]")
		return
	}

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPass := cfg.Database.Pass
	dbName := cfg.Database.Name
	dbSsl := cfg.Database.Ssl
	dbTimeout := cfg.Database.Timeout

	if dbSsl == "" {
		dbSsl = "disable"
	}

	if dbTimeout <= 0 {
		dbTimeout = 5
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d", dbHost, dbPort, dbUser, dbName, dbPass, dbSsl, dbTimeout)
	conn, err := ConnectDB("postgre", dsn, debug)
	if err != nil {
		log.Fatal(err)
	}

	queueHost := cfg.Queue.Host
	queuePort := cfg.Queue.Port
	queueName := cfg.Queue.Name
	dsn = fmt.Sprintf("mongodb://%s:%d", queueHost, queuePort)
	queue, err := ConnectQueue(dsn, c)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = queue.Disconnect(c); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	expired := cfg.Queue.Expired
	if expired < 1 {
		expired = 24
	}
	expiredQueue := time.Duration(expired)

	phoneID := cfg.Chatbot.PhoneID
	urlSendMsg := cfg.Chatbot.Host
	accessToken := cfg.Chatbot.AccessToken
	urlPush := cfg.Chatbot.HostPush
	accountID := cfg.Chatbot.AccountID
	divisionID := cfg.Chatbot.DivisionID
	accessTokenPush := cfg.Chatbot.AccessTokenPush
	downloadFolder := cfg.DownloadFolder

	model := postgre.NewPostgreRepository(c, conn)
	ordersQueue := mongo.NewmongoRepository(c, queue.Database(queueName), "chatbot", expiredQueue)
	chatUcase := usecase.NewChatUcase(model, urlPush, urlSendMsg, phoneID, accessToken, accountID, divisionID, accessTokenPush, downloadFolder, ordersQueue, rdb)

	e := echo.New()
	_HttpDelivery.NewOrderHandler(e, chatUcase, debug)

	go func() {
		if err := e.Start(server); err != nil {
			log.Fatalf("[run_webhook] %s", err)
		}
	}()

	// Wait for the context to be cancelled (e.g., by receiving SIGINT)
	<-c.Done()

	log.Println("Shutting down gracefully...")

	// Create a context with a timeout for shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Shutdown the server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped.")

	return nil
}
