package app

import (
	"context"
	"log"

	"github.com/solumD/auth-test-task/internal/client/db"
	"github.com/solumD/auth-test-task/internal/client/db/pg"
	"github.com/solumD/auth-test-task/internal/closer"
	"github.com/solumD/auth-test-task/internal/config"
	"github.com/solumD/auth-test-task/internal/handler"
	"github.com/solumD/auth-test-task/internal/repository"
	authRepo "github.com/solumD/auth-test-task/internal/repository/auth"
	"github.com/solumD/auth-test-task/internal/service"
	authSrv "github.com/solumD/auth-test-task/internal/service/auth"
	emailSrv "github.com/solumD/auth-test-task/internal/service/email"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	serverConfig config.ServerConfig
	loggerConfig config.LoggerConfig

	dbClient db.Client

	authRepository repository.AuthRepository
	authService    service.AuthService

	emailService service.EmailService

	handler *handler.Handler
}

// NewServiceProvider returns new object of service provider
func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig initializes Postgres config if it was not initialized yet and returns it
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// LoggerConfig initializes logger config if it was not initialized yet and returns it
func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			log.Fatalf("failed to get logger config:%v", err)
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

// ServerConfig initializes http-server config if it was not initialized yet and returns it
func (s *serviceProvider) ServerConfig() config.ServerConfig {
	if s.serverConfig == nil {
		cfg, err := config.NewServerConfig()
		if err != nil {
			log.Fatalf("failed to get http config")
		}

		s.serverConfig = cfg
	}

	return s.serverConfig
}

// DBClient initializes database client config if it was not initialized yet and returns it
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create a db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("postgres ping error: %v", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// AuthRepository initializes auth repository if it was not initialized yet and returns it
func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.New(s.DBClient(ctx))
	}

	return s.authRepository
}

// EmailService initializes email service if it was not initialized yet and returns it
func (s *serviceProvider) EmailService() service.EmailService {
	if s.emailService == nil {
		s.emailService = emailSrv.New()
	}

	return s.emailService
}

// AuthService initializes auth service if it was not initialized yet and returns it
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authSrv.New(s.AuthRepository(ctx), s.EmailService())
	}

	return s.authService
}

// Handler initializes handler if it was not initialized yet and returns it
func (s *serviceProvider) Handler(ctx context.Context) *handler.Handler {
	if s.handler == nil {
		s.handler = handler.New(s.AuthService(ctx))
	}

	return s.handler
}
