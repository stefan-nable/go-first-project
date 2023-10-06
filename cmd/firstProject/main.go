package main

import (
	"Main/internal/http/service"
)

func main() {
	s := service.NewService()
	s.StartWebService()
}
