package service

import (
	"Main/internal/db"
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

func cleanup(dbRef *sql.DB) {
	logger.Printf("Shutting down gracefully and doing cleanup... (flushing database)")

	if dbRef != nil {
		db.FlushDB(dbRef)
		err := dbRef.Close()
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
		var dbRef *sql.DB
		for i := 0; i < 10; i++ {
			select {
			case <-c:
				cleanup(dbRef)
			case dbRef = <-dbChannel:
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

	storage := db.ConnectToDB()
	dbChannel <- storage

	apiManager := routes.NewRouter(storage)
	apiManager.RegisterRoutes(ws)

	logger.Printf("Started api service on port 8080")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
