package application

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"os/signal"
	"subd/application/common/middlewares"
	handler1 "subd/application/forum/delivery"
	repos1 "subd/application/forum/repository"
	useCase1 "subd/application/forum/usecase"
	handler2 "subd/application/user/delivery"
	repos2 "subd/application/user/repository"
	useCase2 "subd/application/user/usecase"
	handler3 "subd/application/thread/delivery"
	repos3 "subd/application/thread/repository"
	useCase3 "subd/application/thread/usecase"
	handler4 "subd/application/service/delivery"
	repos4 "subd/application/service/repository"
	useCase4 "subd/application/service/usecase"
	handler5 "subd/application/post/delivery"
	repos5 "subd/application/post/repository"
	useCase5 "subd/application/post/usecase"
	"syscall"
	"time"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	MaxConn  int    `yaml:"max_conn"`
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
	db       *pgx.ConnPool
}

func initDatabase(config Config) (*pgx.ConnPool, error) {
	pgxConf := pgx.ConnConfig{User: config.DB.User, Password: config.DB.Password,
		Host: config.DB.Host, Port: config.DB.Port, Database: config.DB.Name}
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgxConf, MaxConnections: config.DB.MaxConn})
	if err != nil {
		return nil, fmt.Errorf("error in connection pool create: %s", err)
	}
	return conn, nil
}

func NewApp(config Config) *App {

	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := Logger{Info: infoLogger, Error: errorLogger}
	infoLogger.SetLogLevel(logger.DebugLevel)

	conn, err := initDatabase(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %s\n", err)
		os.Exit(1)
	}

	return &App{config: config, log: log, doneChan: make(chan bool, 1), db: conn}
}

func initRouting(a *App) *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/health", healthCheck)
	forumRepos := repos1.NewPgRepository(a.db)
	forumCase := useCase1.NewForumUseCase(forumRepos)
	handler1.NewForumHandler(router, forumCase)

	userRepos := repos2.NewPgRepository(a.db)
	userCase := useCase2.NewUserUseCase(userRepos)
	handler2.NewUserHandler(router, userCase)

	threadRepos := repos3.NewPgRepository(a.db)
	threadCase := useCase3.NewThreadUseCase(threadRepos)
	handler3.NewThreadHandler(router, threadCase)

	serviceRepos := repos4.NewPgRepository(a.db)
	serviceCase := useCase4.NewServiceUseCase(serviceRepos)
	handler4.NewUserHandler(router, serviceCase)

	postRepos := repos5.NewPgRepository(a.db, userRepos, forumRepos, threadRepos)
	postCase := useCase5.NewPostUseCase(postRepos)
	handler5.NewPostHandler(router, postCase)
	return router
}

func (a *App) Run() {
	router := initRouting(a)
	go func() {
		a.log.Info.Infof("Start listening on %s", a.config.Listen)
		if err := fasthttp.ListenAndServe(a.config.Listen, middlewares.JsonRequestHandler(router.Handler)); err != nil && err != http.ErrServerClosed {
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
	defer a.db.Close()
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
