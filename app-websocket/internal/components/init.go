package components

import (
	"app-websocket/clients/sso"
	"app-websocket/internal/broker/kafka"
	"app-websocket/internal/config"
	"app-websocket/internal/ports"
	authhandle "app-websocket/internal/ports/http/auth"
	chathandle "app-websocket/internal/ports/http/chat"
	"app-websocket/internal/ports/ws"
	"app-websocket/internal/services/auth"
	"app-websocket/internal/services/client"
	"app-websocket/internal/services/room"
	clientstor "app-websocket/internal/storage/client"
	"app-websocket/internal/storage/connector"
	roomstor "app-websocket/internal/storage/room"

	"github.com/sirupsen/logrus"
)

type Components struct {
	HTTPServer    *ports.Server
	Database      *connector.Database
	KafkaProducer *kafka.Producer
	KafkaConsumer *kafka.Consumer
}

func InitComponents() (*Components, error) {
	if err := config.Load(); err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	pgConf, err := config.NewDBConfig()
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	pg, err := connector.NewDatabase(pgConf)
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}

	kfConf, err := config.NewKafkaConfig()
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	kafkaProducer, err := kafka.NewProducer(*kfConf)
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	kafkaConsumer, err := kafka.NewConsumer(*kfConf)
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}

	ssocfg, err := config.NewSSOConfig()
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	ssogrpcServiceClient, err := sso.NewSSOServiceClient(logrus.New(), *ssocfg)
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}

	hub := ws.NewHub(kafkaConsumer)

	authService := auth.NewAuthService(ssogrpcServiceClient)

	clientRepo := clientstor.NewClientRepository(pg)
	clientService := client.NewClientService(clientRepo, kafkaProducer, hub)

	roomRepo := roomstor.NewRoomRepository(pg)
	roomService := room.NewRoomService(roomRepo)

	authHandler := authhandle.NewHandler(authService)
	chatHandler := chathandle.NewHandler(roomService, clientService)

	serverCfg, err := config.NewHTTPServerConfig()
	if err != nil {
		logrus.Error("initComponents", err)
		return nil, err
	}
	httpServer := ports.NewServer(serverCfg, *authHandler, *chatHandler, hub)
	return &Components{
		HTTPServer:    httpServer,
		Database:      pg,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
	}, nil

}
