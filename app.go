package main

import (
	"./models"
	"./storager"
	"github.com/insionng/torgo"
)

func main() {
	models.CreateDb()
	torgo.SetStaticPath("/static", "./static")

	torgo.Router(`/getfile/:filename`, &storager.RStorageHandler{})

	torgo.Router(`/setfile/:filename`, &storager.WStorageHandler{})

	torgo.Router("/", &storager.UploaderHandler{})
	torgo.Router("/uploader", &storager.UploaderHandler{})

	torgo.Run()
}
