package ports

import (
	"app-websocket/internal/config"
	authhandle "app-websocket/internal/ports/http/auth"
	chathandle "app-websocket/internal/ports/http/chat"
	"app-websocket/internal/ports/ws"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
	hub    *ws.Hub
}

func NewServer(cfg config.HTTPServerConfig, authHandler authhandle.Handler, chatHandler chathandle.Handler, hub *ws.Hub) *Server {
	return &Server{
		server: &http.Server{
			Addr:           cfg.Address,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        InitRoutes(authHandler, chatHandler),
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
		},
		hub: hub,
	}
}

// TODO: implement logic with running ws.Hub
func (s *Server) Run(ctx context.Context) error {
	errResult := make(chan error)
	go func() {
		logrus.Info(fmt.Sprintf("starting listening: %s", s.server.Addr))

		errResult <- s.server.ListenAndServe()
	}()

	go func() {
		s.hub.Run(ctx)
	}()

	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()

	case err = <-errResult:
	}
	return err
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func InitRoutes(authHandler authhandle.Handler, chatHandler chathandle.Handler) *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}
	ws := router.Group("/ws", authHandler.UserIdentity)
	{
		ws.POST("/createRoom", chatHandler.CreateRoom)
		ws.GET("/joinRoom/:roomID", chatHandler.JoinRoom)
		ws.GET("/getRooms", chatHandler.JoinRoom)
	}
	return router
}
