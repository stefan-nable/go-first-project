package model

import (
	"fmt"
	"time"
)

type Log struct {
	Message   string
	Timestamp time.Time
}

func (l Log) PrintLog() {
	fmt.Println(l.Timestamp.Format("2006-01-02 15:04:05") + " " + l.Message)
}
