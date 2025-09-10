package main

import (
	"os"
	"os/signal"
	"subscription/internal/config"
	"subscription/internal/logger"
	"subscription/internal/repository"
	"subscription/internal/server"
	"subscription/internal/service"
	"syscall"
)

func waitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
func main() {
	cfg := config.MustNew()

	lg := logger.MustNew(&cfg.Log)
	defer lg.Sync()

	repo := repository.MustNew(lg, &cfg.Db)
	defer repo.Close()

	svc := service.New(repo, lg)

	srv := server.New(svc, lg, &cfg.Srv)
	srv.Start()
	defer srv.Stop()

	waitSignal()

}
