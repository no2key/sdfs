package main

import (
	"./models"
	"./server"
	"github.com/insionng/torgo"
)

func main() {
	models.CreateDb()
	torgo.Router(`/getfile/:filename`, &server.RStorageHandler{})
	torgo.Router(`/setfile/:filename`, &server.WStorageHandler{})
	torgo.Run()
}
