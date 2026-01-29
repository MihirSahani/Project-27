package main

import (
	"net/http"
	"time"

	"github.com/MihirSahani/Project-27/internal"
	"github.com/MihirSahani/Project-27/internal/jwt"
	"github.com/MihirSahani/Project-27/storage"
	"go.uber.org/zap"
)

type LogCategory int

const (
	ErrorLog LogCategory = iota
	WarnLog
	InfoLog
)

type Application struct {
	config AppConfig
	logger *zap.Logger
	authenticator internal.Authenticator
	storageManager *storage.StorageManager
}

func NewApp() *Application {
	return &Application{
		config: LoadServerConfig(),
		logger: NewLogger(),
		authenticator: jwt.NewJWTAuthenticator(),
		storageManager: storage.NewStorageManager(),
	}
}

func (app *Application) Run() {
	router := app.mount()
	server := &http.Server{
		Addr:         app.config.ServerAddress,
		Handler:      router,
		IdleTimeout:  time.Duration(app.config.IdleTimeout),
		ReadTimeout:  time.Duration(app.config.ReadTimeout),
		WriteTimeout: time.Duration(app.config.WriteTimeout),
	}
	app.logger.Info("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (app *Application) ErrorLogger(message string, err error, status int, w http.ResponseWriter, category LogCategory) {
	switch category {
	case ErrorLog:
		app.logger.Error(message, zap.Error(err))
	case WarnLog:
		app.logger.Warn(message, zap.Error(err))
	default:
		app.logger.Info(message, zap.Error(err))
	}
	app.errorJSON(w, status)
}

func (app *Application) Close() {
	app.logger.Info("Closing server")
	app.logger.Sync()

	app.logger.Info("Closing Database connection")
	app.storageManager.Close()
}
