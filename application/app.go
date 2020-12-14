package application

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Listen string   `yaml:"listen"`
	DB     DBConfig `yaml:"db"`
}

type Logger struct {
	Info  *logger.Logger
	Error *logger.Logger
}

type App struct {
	config   Config
	log      Logger
	doneChan chan bool
	db       *pgx.Conn
}

func NewApp(config Config) *App {

	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := Logger{Info: infoLogger, Error: errorLogger}
	infoLogger.SetLogLevel(logger.DebugLevel)

	conn, err := pgx.Connect(pgx.ConnConfig{User: config.DB.User, Password: config.DB.Password,
		Host: config.DB.Host, Port: config.DB.Port, Database: config.DB.Name})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	//router.GET("/health", healthCheck())
	return &App{config: config, log: log, doneChan: make(chan bool, 1), db: conn}
}

func (a *App) Run() {
	handlersPool := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/forum/create":
		case "/health":
			a.log.Info.Notice("Get connection-test request...")
			healthCheck(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	go func() {
		a.log.Info.Infof("Start listening on %s", a.config.Listen)
		if err := fasthttp.ListenAndServe(a.config.Listen, handlersPool); err != nil && err != http.ErrServerClosed {
			a.log.Error.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-a.doneChan:
	}
	a.log.Info.Info("Shutdown Server (timeout of 1 seconds) ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	<-ctx.Done()

	a.log.Info.Info("Server exiting")
}

func (a *App) Close() {
	a.doneChan <- true
}

func healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}
