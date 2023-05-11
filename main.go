package main

import (
	"github.com/alfaa19/gin-stegano/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//router
	r.POST("/stegano/encode", controller.Encode)
	r.POST("/stegano/decode", controller.Decode)
	r.Run(":8001")
}
