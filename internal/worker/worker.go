package worker

import (
	"Main/internal/model"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"
)

type W interface {
	StartWorker(wg *sync.WaitGroup, inputChannel chan *InputData, resultChannel chan model.Log)
}

type Worker struct{}

func (w *Worker) StartWorker(wg *sync.WaitGroup, inputChannel chan *InputData, resultChannel chan model.Log) {
	defer wg.Done()
	var res []string
	inputData := <-inputChannel

	logger.Print("Started worker " + strconv.Itoa(inputData.WorkerNumber))

	file := getWorkerFile(inputData.InputFile, inputData.NumWorkers, inputData.WorkerNumber)

	for _, line := range file {
		a, err := strconv.Atoi(strings.Split(line, " ")[0])
		if err != nil {
			resultChannel <- model.Log{
				Message: "Error converting string to int: " + err.Error(),
			}
			return
		}
		b, err := strconv.Atoi(strings.Split(line, " ")[1])
		if err != nil {
			resultChannel <- model.Log{
				Message: "Error converting string to int: " + err.Error(),
			}
			return
		}

		res = append(res, strconv.FormatFloat(inputData.MathFunc(float64(a), float64(b)), 'f', 1, 64))
	}

	ans := strings.Join(res, ", ")

	log := model.Log{
		Message:   "I'm a worker and I computed: " + ans,
		Timestamp: time.Now(),
	}

	resultChannel <- log
}

func getWorkerFile(file *[]string, numWorkers int, workerNumber int) []string {
	totalLines := len(*file)
	linesPerWorker := totalLines / numWorkers

	startingLine := workerNumber * linesPerWorker
	var endingLine int

	if workerNumber == numWorkers-1 {
		endingLine = totalLines
	} else {
		endingLine = (workerNumber + 1) * linesPerWorker
	}

	ans := (*file)[startingLine:endingLine]

	return ans
}
