package main

import (
	"github.com/gin-gonic/gin"
	//"net/http"
	//"io"
    //"fmt"
)

func web(){
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "web.html", nil)
	 })

	router.POST("/demo3", func(c *gin.Context) {
		user := Config{} 
		err := c.ShouldBind(&user)  
		
		if err != nil {
		   c.String(500, "Error")
		} else {
		   c.String(200, "ok")
		}
	 })
	router.Run(":8080")
}