package main

import (
	"Main/internal/model"
	"Main/internal/utils"
	"Main/internal/worker/bigBossWorker"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
)

func main() {
	resultChannel := make(chan model.Log)
	wg := new(sync.WaitGroup)

	content, err := os.ReadFile("test/input")
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}
	file := strings.Split(string(content), "\n")

	numWorkers, err := utils.GetNumWorkers(len(file))
	if err != nil {
		fmt.Println("Error getting number of workers:", err)
		return
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		if rand.Int()%2 == 0 {
			go bigBossWorker.StartWorker(wg, resultChannel, utils.GetWorkerFile(&file, numWorkers, i), math.Min)
		} else {
			go bigBossWorker.StartWorker(wg, resultChannel, utils.GetWorkerFile(&file, numWorkers, i), math.Max)
		}
	}

	utils.SaveLogsToDB(numWorkers, resultChannel)

	wg.Wait()
	close(resultChannel)
}
