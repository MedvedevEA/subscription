package server

import (
	"context"
	"subscription/internal/config"
	"subscription/internal/service"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Server struct {
	app      *fiber.App
	bindAddr string
	lg       *zap.SugaredLogger
}

func New(svc service.Service, lg *zap.SugaredLogger, cfg *config.Srv) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		WriteTimeout: cfg.WriteTimeout,
	})
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: lg.Desugar(),
	}))

	appGroup := app.Group("/api/v1")

	appGroup.Post("/services", svc.AddService)
	appGroup.Get("/services/:id", svc.GetService)
	appGroup.Get("/services", svc.GetServices)
	appGroup.Put("/services/:id", svc.UpdateService)
	appGroup.Delete("/services/:id", svc.RemoveService)

	appGroup.Post("/subscriptions", svc.AddSubscription)
	appGroup.Get("/subscriptions/:id", svc.GetSubscription)
	appGroup.Get("/subscriptions", svc.GetSubscriptions)
	appGroup.Post("/subscriptions/total", svc.GetSubscriptionTotal)
	appGroup.Put("/subscriptions/:id", svc.UpdateSubscription)
	appGroup.Delete("/subscriptions/:id", svc.RemoveSubscription)

	return &Server{
		app:      app,
		bindAddr: cfg.Addr,
		lg:       lg,
	}

}

func (s *Server) Start() {
	s.lg.Infof("server start (bind address: %s)", s.bindAddr)
	go func() {
		if err := s.app.Listen(s.bindAddr); err != nil {
			s.lg.Errorf("failed to listen server: %v", err)
		}
	}()

}
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		s.lg.Errorf("failed to shutdown server: %v", err)
		return
	}
	s.lg.Info("server shutdown")

}
