package routes

import (
	"Main/internal/model"
	"Main/internal/worker"
	"encoding/json"
	"github.com/emicklei/go-restful"
	logger "github.com/sirupsen/logrus"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type numberOfWorkersStruct struct {
	NumberOfWorkers int `json:"numberOfWorkers"`
}

func (r *Router) PostWorkersHandler(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		logger.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusBadRequest, restful.NewError(http.StatusBadRequest, "Bad request\n"))
		return

	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		logger.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "Internal server error\n"))
		return
	}

	var numberOfWorkers = numberOfWorkersStruct{-1}
	if err := json.Unmarshal(data, &numberOfWorkers); err != nil || numberOfWorkers.NumberOfWorkers == -1 {
		logger.Printf("[ERROR] Couldn't unmarshal request body")
		resp.WriteError(http.StatusBadRequest, restful.NewError(http.StatusBadRequest, "Bad request\n"))
		return
	} else if numberOfWorkers.NumberOfWorkers == 0 {
		logger.Printf("[ERROR] Number of workers can't be 0")
		resp.WriteError(http.StatusBadRequest, restful.NewError(http.StatusBadRequest, "Number of workers must be higher than 0\n"))
		return
	}

	resp.WriteAsJson("Processing started with " + strconv.Itoa(numberOfWorkers.NumberOfWorkers) + " workers.")

	r.startProcessingOnPostHandler(numberOfWorkers.NumberOfWorkers)
}

func (r *Router) startProcessingOnPostHandler(numWorkers int) {
	resultChannel := make(chan model.Log)
	wg := new(sync.WaitGroup)

	logger.Printf("Started processing with %d workers", numWorkers)

	content, err := os.ReadFile("test/input")
	if err != nil {
		logger.Error("Error reading file: ", err)
		return
	}
	file := strings.Split(string(content), "\n")

	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		workerInputData := worker.InputData{
			WorkerNumber: i,
			NumWorkers:   numWorkers,
			InputFile:    &file,
		}
		
		if rand.Int()%2 == 0 {
			workerInputData.MathFunc = math.Min
		} else {
			workerInputData.MathFunc = math.Max
		}

		go worker.StartWorker(wg, resultChannel, &workerInputData)
	}

	//memory.FlushDB(r.db)

	for i := 0; i < numWorkers; i++ {
		log := <-resultChannel
		model.Log.PrintLog(log)

		_, err := r.db.Exec("INSERT INTO log (timestamp, message) VALUES (?, ?)",
			log.Timestamp.Format("2006-01-02 15:04:05"),
			log.Message)

		if err != nil {
			panic(err.Error())
		}
	}

	wg.Wait()
	close(resultChannel)
}
