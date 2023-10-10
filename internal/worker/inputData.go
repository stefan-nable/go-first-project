package worker

type InputData struct {
	WorkerNumber int
	NumWorkers   int
	MathFunc     func(float64, float64) float64
	InputFile    *[]string
}
