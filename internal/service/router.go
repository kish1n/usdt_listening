package service

import (
	"github.com/go-chi/chi"
	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/data/pg"
	"github.com/kish1n/usdt_listening/internal/service/handlers"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"net/http"
	"sync"
)

func (s *service) router(cfg config.Config) (chi.Router, error) {
	r := chi.NewRouter()
	logger := cfg.Log()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxDB(pg.NewMasterQ(cfg.DB())),
			helpers.CtxServiceConfig(cfg.ServiceConfig()),
		),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handlers.ListenTransfers(cfg)
	}()

	r.Route("/", func(r chi.Router) {
		r.Get("/from/{sender}", handlers.SortBySender)
		r.Get("/to/{recipient}", handlers.SortByRecipient)
		r.Get("/by/{address}", handlers.SortByAddress)
	})

	logger.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		return r, err
	}

	return r, nil
}
