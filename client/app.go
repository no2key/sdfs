package main

import (
	"./client"
	//"./models"
	"github.com/insionng/torgo"
)

func main() {
	//models.CreateDb()
	torgo.SetStaticPath("/static", "./static")

	torgo.Router("/", &client.UploaderHandler{})
	torgo.Router("/uploader", &client.UploaderHandler{})

	torgo.Run()
}
