package bigBossWorker

import (
	"Main/internal/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StartWorker(wg *sync.WaitGroup, resultChannel chan model.Log, file []string, mathFunc func(float64, float64) float64) {
	defer wg.Done()
	var res []string

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

		res = append(res, strconv.FormatFloat(mathFunc(float64(a), float64(b)), 'f', 1, 64))
	}

	ans := strings.Join(res, ", ")

	log := model.Log{
		Message:   "I'm a worker and I computed: " + ans,
		Timestamp: time.Now(),
	}

	resultChannel <- log
}
