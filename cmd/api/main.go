package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/api"
	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/store"
)

func main() {
	logger := log.New(os.Stdout, "ai-planner ", log.LstdFlags)

	st, err := store.NewWithAutoInit()
	if err != nil {
		logger.Fatalf("connect to MariaDB: %v", err)
	}
	defer st.Close()

	addr := os.Getenv("AI_PLANNER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      api.New(st, logger).Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	logger.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("shutdown error: %v", err)
	}
}
