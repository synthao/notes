package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/synthao/notes/internal/adapter/mysql/note/repository"
	"github.com/synthao/notes/internal/config"
	port "github.com/synthao/notes/internal/port/http/note"
	"github.com/synthao/notes/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func main() {
	fx.New(
		fx.Provide(
			fiber.New,
			config.NewServerConfig,
			config.NewDBConfig,
			config.NewLoggerConfig,
			newLogger,
			repository.NewRepository,
			service.NewService,
			port.NewHandler,
			newDBConnection,
		),
		fx.Invoke(createHTTPServer),
	).Run()
}

func newDBConnection(cnf *config.DB) *sqlx.DB {
	return sqlx.MustConnect("mysql", cnf.GetMysqlDSN())
}

func newLogger(cnf *config.Logger) (*zap.Logger, error) {
	atomicLogLevel, err := zap.ParseAtomicLevel(cnf.Level)
	if err != nil {
		return nil, err
	}

	atom := zap.NewAtomicLevelAt(atomicLogLevel.Level())
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
	), nil
}

func createHTTPServer(lc fx.Lifecycle, app *fiber.App, handler *port.Handler, cnf *config.Server) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	handler.InitRoutes()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(fmt.Sprintf(":%d", cnf.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("==> stopping server ...")

			return app.Shutdown()
		},
	})
}
