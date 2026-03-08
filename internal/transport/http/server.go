package http

import (
	"fmt"

	"github.com/Nikkoz/task-service/internal/config"
	"github.com/Nikkoz/task-service/internal/transport/http/auth"
	"github.com/Nikkoz/task-service/internal/transport/http/task"
	"github.com/gin-gonic/gin"
)

type (
	Server struct {
		router *gin.Engine

		options Options
	}

	Options struct {
		Notify chan error
	}
)

func NewServer(taskService task.Service, authService auth.Service, isProd bool, authCfg config.Auth, o Options) *Server {
	taskHandler := task.NewHandler(taskService)
	authHandler := auth.NewHandler(authService)

	route := newRouter(taskHandler, authHandler, isProd, authCfg)

	s := &Server{
		router: route,
	}

	s.setOptions(o)

	return s
}

func (s *Server) setOptions(options Options) {
	if options.Notify == nil {
		s.options.Notify = make(chan error, 1)
	}

	if s.options != options {
		s.options = options
	}
}

func (s *Server) Run(cfg config.Http) {
	go func() {
		defer close(s.options.Notify)

		s.options.Notify <- s.router.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	}()
}

func (s *Server) Notify() <-chan error {
	return s.options.Notify
}
