package main

import (
    "context"
    stdhttp "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/diegodesousas/go-devkit/pkg/gen"
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/go-devkit/pkg/log"
    "github.com/diegodesousas/go-devkit/pkg/metrics"
    "github.com/diegodesousas/greeter/internal/infra/clock"
    "github.com/diegodesousas/greeter/internal/infra/database"
    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/greeter/internal/infra/http/routes"
    "github.com/spf13/viper"
)

func bootstrapConfig() error {
    os.Setenv("TZ", "UTC")

    viper.SetConfigType("env")
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return err
    }

    return nil
}

func bootstrapLogger() log.Logger {
    levelMap := map[string]log.Level{
        "debug":   log.DebugLevel,
        "warning": log.WarnLevel,
        "info":    log.InfoLevel,
        "error":   log.ErrorLevel,
    }

    level, ok := levelMap[viper.GetString("LOG_LEVEL")]
    if !ok {
        level = log.InfoLevel
    }

    return log.New(
        log.WithLevel(level),
        log.WithJsonFormat(),
    )
}

func bootstrapServer(routeOpt httpserver.Option, logger log.Logger, metricsClient metrics.Metric) httpserver.Server {
    return httpserver.New(
        httpserver.WithAPM(viper.GetBool("DD_TRACE_APM_ENABLED")),
        httpserver.WithName("greeter"),
        httpserver.WithPort(viper.GetString("HTTP_PORT")),
        httpserver.WithMiddlewares(
            httpserver.Logger(logger),
            metrics.Metrics(metricsClient),
            httpserver.RequestID,
            httpserver.TraceID(gen.UUIDGenerator()),
            httpserver.ContentTypeJson(),
            httpserver.Compress(),
            httpserver.AllowAll(),
        ),
        httpserver.WithErrorHandler(infrahttp.ErrorHandler),
        httpserver.WithHTTPServerReadTimeout(time.Second*60),
        routeOpt,
    )
}

func bootstrapDatabase() (database.Connection, error) {
	conn, err := database.NewPostgresConnection()
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func bootstrapRoutes(repos database.Repositories) httpserver.Option {
    appClock := clock.New()

    var routeList []httpserver.Route
    routeList = append(routeList, routes.Health()...)
    routeList = append(routeList, routes.Greeting(appClock, repos.Greeting)...)

    return httpserver.WithRoutes(routeList...)
}

func main() {
    if err := bootstrapConfig(); err != nil {
        log.Warn(context.Background(), err.Error())
    }

    logger := bootstrapLogger().WithFields(log.Field{
        Key:   "env",
        Value: viper.GetString("ENV"),
    })

    ctx := log.WithLogger(context.Background(), logger)

    statsdClient, err := metrics.New()
    if err != nil {
        log.FatalError(ctx, err)
    }

    conn, err := bootstrapDatabase()
    if err != nil {
        log.FatalError(ctx, err)
    }

    repos := database.NewRepositories(conn)
    server := bootstrapServer(bootstrapRoutes(repos), logger, statsdClient)

    log.Info(ctx, "server starting...")
    shutdown := server.Run()

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := <-server.ShutdownListener(); err != nil && err != stdhttp.ErrServerClosed {
            interrupt <- syscall.SIGTERM
        }
    }()

    log.Info(ctx, "server running")
    <-interrupt

    if err := shutdown(ctx); err != nil {
        log.Fatal(ctx, err.Error())
    }

    log.Info(ctx, "server shutdown completed")
}
