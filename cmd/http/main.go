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
    _ "github.com/diegodesousas/greeter/docs"
    "github.com/diegodesousas/greeter/internal/infra/clock"
    "github.com/diegodesousas/greeter/internal/infra/database"
    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/greeter/internal/infra/http/routes"
    "github.com/diegodesousas/greeter/internal/infra/shutdown"
    "github.com/spf13/viper"
)

// shutdownTimeout bounds graceful shutdown; keep it below the orchestrator's
// grace period (30s by default on Kubernetes) so resources are released before
// the process is killed.
const shutdownTimeout = 15 * time.Second

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
        log.WithJSONFormat(),
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
            httpserver.ContentTypeJSON(),
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

func bootstrapRoutes(conn database.Connection, repos database.Repositories) httpserver.Option {
    appClock := clock.New()

    var routeList []httpserver.Route
    routeList = append(routeList, routes.Health(conn)...)
    routeList = append(routeList, routes.Greeting(appClock, repos.Greeting)...)
    routeList = append(routeList, routes.Docs()...)

    return httpserver.WithRoutes(routeList...)
}

//	@title			Greeter API
//	@version		1.0
//	@description	Serviço de saudações - exemplo de Clean Architecture em Go
//	@host			localhost:3000
//	@BasePath		/
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
    server := bootstrapServer(bootstrapRoutes(conn, repos), logger, statsdClient)

    log.Info(ctx, "server starting...")
    stopServer := server.Run()

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := <-server.ShutdownListener(); err != nil && err != stdhttp.ErrServerClosed {
            interrupt <- syscall.SIGTERM
        }
    }()

    log.Info(ctx, "server running")
    <-interrupt

    err = shutdown.Graceful(ctx, shutdownTimeout,
        shutdown.Step(stopServer),
        func(context.Context) error { return conn.Close() },
    )
    if err != nil {
        log.Error(ctx, err)
        os.Exit(1)
    }

    log.Info(ctx, "server shutdown completed")
}
