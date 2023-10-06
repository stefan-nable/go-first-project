package service

import (
	"Main/internal/http/memory"
	"Main/internal/http/routes"
	"database/sql"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/emicklei/go-restful"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func cleanup(db *sql.DB) {
	logger.Printf("Shutting down gracefully and doing cleanup... (flushing database)")

	if db != nil {
		memory.FlushDB(db)
		err := db.Close()
		if err != nil {
			return
		}
	}
	os.Exit(1)
}

func gracefulShutdown(dbChannel chan *sql.DB) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(dbChannel chan *sql.DB) {
		var db *sql.DB
		for i := 0; i < 10; i++ {
			select {
			case <-c:
				cleanup(db)
			case db = <-dbChannel:
			}
		}
	}(dbChannel)
}

func loggerSetup() {
	formatter := runtime.Formatter{ChildFormatter: &logger.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}}
	formatter.Line = true
	logger.SetFormatter(&formatter)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.InfoLevel)
}

func (s *Service) StartWebService() {
	dbChannel := make(chan *sql.DB)
	gracefulShutdown(dbChannel)
	loggerSetup()

	ws := new(restful.WebService)
	restful.Add(ws)

	storage := memory.ConnectToDB()
	dbChannel <- storage

	apiManager := routes.NewRouter(storage)
	apiManager.RegisterRoutes(ws)

	logger.Printf("Started api service on port 8080")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
