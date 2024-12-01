package main

import (
	"context"
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/services"
	"costumers-api/infrastructure"
	"costumers-api/infrastructure/databases/costumers_db"
	"costumers-api/infrastructure/repositories"
	httpapi "costumers-api/presentation/http_api"
	costumercontroller "costumers-api/presentation/http_api/costumer_controller"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CostumersServiceApp struct {
	stopSigChan           chan os.Signal
	costumersDbConnection *sql.DB
	rabbitMqConnection    *amqp.Connection
	unitOfWorkFactory     abstractions.IUnitOfWorkFactory
	messagePublisher      dependencies.IMessagePublisher
	costumersRepository   dependencies.ICostumerRepository
	searchCostumersQuery  dependencies.ISearchCostumerQuery
	eventRepository       dependencies.IEventRepository
	cacheInstance         dependencies.IDistributedCacheRepository
	eventHandlerWorker    services.EventHandler
	costumerController    costumercontroller.CostumerController
	healthController      httpapi.HealthController
}

func NewCostumersServiceApp(stopSigChan chan os.Signal) (CostumersServiceApp, error) {

	app := CostumersServiceApp{
		stopSigChan: stopSigChan,
	}

	costumerDbErr := app.connectCostumersDb()
	rmqErr := app.connectRabbitMq()

	err := errors.Join(costumerDbErr, rmqErr)

	if err != nil {
		return CostumersServiceApp{}, err
	}

	app.initInfrastructure()
	app.initServices()
	app.initControllers()
	return app, nil
}

func (app *CostumersServiceApp) connectCostumersDb() error {
	costumersDbConnection, costumersDbErr := costumers_db.ConnectPostgresDB()

	if costumersDbErr != nil {
		fmt.Println("Failed to connect with database")
		return costumersDbErr
	}

	app.costumersDbConnection = costumersDbConnection
	return costumersDbErr
}

func (app *CostumersServiceApp) connectRabbitMq() error {
	rabbitConnectionString := os.Getenv("RABBIT_MQ_CONNECTION")
	rmqConn, err := amqp.Dial(rabbitConnectionString)
	if err != nil {
		fmt.Println("Failed to connect with rabbitmq")
	}
	app.rabbitMqConnection = rmqConn
	return err
}

func (app *CostumersServiceApp) initInfrastructure() {

	unitOfWorkFactory := infrastructure.NewUnitOfWorkFactory(app.costumersDbConnection)
	messagePublisher := infrastructure.NewMessagePublisher(app.rabbitMqConnection)

	costumersRepository := repositories.CreateCostumerRepository(app.costumersDbConnection, context.Background())
	searchCostumersQuery := repositories.CreateSearchCostumerQuery(app.costumersDbConnection, context.Background())
	eventRepository := repositories.NewEventRepository(app.costumersDbConnection, context.Background())
	cacheInstance := repositories.NewRedisCache(context.Background())

	app.unitOfWorkFactory = unitOfWorkFactory
	app.messagePublisher = messagePublisher
	app.costumersRepository = costumersRepository
	app.searchCostumersQuery = searchCostumersQuery
	app.eventRepository = eventRepository
	app.cacheInstance = cacheInstance
}

func (app *CostumersServiceApp) initServices() {
	eventHandlerWorker := services.NewEventHandler(app.eventRepository, app.messagePublisher)
	app.eventHandlerWorker = eventHandlerWorker
}

func (app *CostumersServiceApp) initControllers() {
	costumerController := costumercontroller.NewCostumerController(app.costumersRepository, app.searchCostumersQuery, app.cacheInstance, app.unitOfWorkFactory, app.eventRepository)
	healthController := httpapi.NewHealthController()

	app.costumerController = costumerController
	app.healthController = healthController
}

func (app *CostumersServiceApp) StartWorkers() {
	go func() {
		app.eventHandlerWorker.RunJob(app.stopSigChan)
	}()
}

func (app *CostumersServiceApp) RegisterControllers(mux *mux.Router) {
	app.costumerController.Inject(mux)
	app.healthController.Inject(mux)
}

func (app *CostumersServiceApp) Close() {
	app.costumersDbConnection.Close()
	app.rabbitMqConnection.Close()
}
