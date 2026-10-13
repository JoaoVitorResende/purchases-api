package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"purchases-api/internal/application/query"
	"purchases-api/internal/application/command"
	"purchases-api/internal/infrastructure/postgres"
	"purchases-api/internal/interfaces/httpapi"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL nao definida")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	// Infrastructure
	purchaseRepo := postgres.NewPurchaseRepository(pool)
	purchaseReader := postgres.NewPurchaseReader(pool)
	catalogStore := postgres.NewCatalogStore(pool)

	// Application
	createPurchase := command.NewCreatePurchaseHandler(purchaseRepo, catalogStore, time.Now)
	getPurchase := query.NewGetPurchaseHandler(purchaseReader)

	// Interfaces
	handler := httpapi.NewPurchaseHandler(createPurchase, getPurchase)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewRouter(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("purchases-api na porta %s", port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
