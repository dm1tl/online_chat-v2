package ports

import (
	"app-websocket/internal/config"
	"app-websocket/internal/ports/http/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	//ws.Hub
}

func NewServer(cfg config.HTTPServerConfig, authHandler auth.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           cfg.Address,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        InitRoutes(authHandler),
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
		},
	}
}

// TODO: implement logic with running ws.Hub
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func InitRoutes(authHandler auth.Handler) *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}
	return router
}
