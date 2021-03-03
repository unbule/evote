package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kongyixueyuan.com/evote/web/controllers"
)

func main()  {

	//fmt.Printf("main.....",con)
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	registerRouter(r)
	r.Run(":8000")
}

func registerRouter(router *gin.Engine){
	new(controllers.Controller).Router(router)
}


