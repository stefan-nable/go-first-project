package utils

import "fmt"

func GetNumWorkers(fileLength int) (int, error) {
	fmt.Println("Please input the number of workers:")

	var numWorkers int
	_, err := fmt.Scanln(&numWorkers)
	if err != nil {
		fmt.Println("Error reading number of workers:", err)
		return 0, err
	} else if numWorkers > fileLength {
		numWorkers = fileLength
	}

	fmt.Println("The number of workers is:", numWorkers)

	return numWorkers, nil
}

func GetWorkerFile(file *[]string, numWorkers int, workerNumber int) []string {
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
