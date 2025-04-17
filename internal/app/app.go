package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/solumD/auth-test-task/internal/closer"
	"github.com/solumD/auth-test-task/internal/config"
	"github.com/solumD/auth-test-task/internal/logger"

	"github.com/go-chi/chi/v5"
)

const configPath = ".env"

// App object of an app
type App struct {
	serviceProvider *serviceProvider
	server          *http.Server
}

// NewApp returns new App object
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts an App
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	closer.Add(a.shutdownServer)
	if err := a.runServer(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	a.initServiceProvider()
	logger.Init(logger.GetCore(logger.GetAtomicLevel(a.serviceProvider.LoggerConfig().Level())))
	a.initServer(ctx)

	return nil
}

func (a *App) initConfig() error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider() {
	a.serviceProvider = NewServiceProvider()
}

// initServer inits router and routes of a server
func (a *App) initServer(ctx context.Context) {
	router := chi.NewRouter()

	router.Route("/token", func(r chi.Router) {
		r.Get("/generate", a.serviceProvider.Handler(ctx).GenerateTokens(ctx))
		r.Post("/refresh", a.serviceProvider.Handler(ctx).RefreshTokens(ctx))
	})

	srv := &http.Server{
		Addr:    a.serviceProvider.ServerConfig().Address(),
		Handler: router,
	}

	a.server = srv
}

func (a *App) runServer() error {
	log.Printf("server is running on %s", a.serviceProvider.ServerConfig().Address())

	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	log.Println("server stopped")

	return nil
}

// shutdownServer gracefully shutdowns server
func (a *App) shutdownServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
