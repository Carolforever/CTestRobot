package main

import (
	"log"
	"os"
	//"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	//"github.com/tidwall/sjson"
	"net/http"
	//"io"
    "fmt"
)

func router(){
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	//get路由
	router.GET("/lsc", func(c *gin.Context) {
		c.HTML(200, "initial.html", nil)
	 })

	//本地项目POST路由
	router.POST("/local", func(c *gin.Context) {
		user := Config{} 
		err := c.ShouldBind(&user)  
		
		if err != nil {
		   c.String(500, "Error")
		} else {
			flag.Parse()
			if *flagConfig == "" {
				log.Fatalf("No config file specified")
			}
			config := parseConfig(*flagConfig, user)
			go CheckAll(config)
			c.HTML(200, "initial.html", nil)
		}
	 })

	//上传项目POST路由
	router.POST("/upload", func(c *gin.Context) {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        filePath := robot_dir + "/projects/" + file.Filename

        // 保存上传的文件到指定路径
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.String(http.StatusInternalServerError, fmt.Sprintf("保存文件失败: %s", err.Error()))
            return
        }

		/*
		err, _ = RunCommand(robot_dir + "/projects", "tar", "-zxvf", filePath)
		if err != nil {
			log.Println("Unzip .gz file: failed for :", err)
		}

		err, _ = RunCommand(robot_dir + "/projects", "rm", "-rf", filePath)
		if err != nil {
			log.Println("delete .gz file: failed for :", err)
		}
		
		c.HTML(200, "initial.html", nil)
		*/
    })	

	router.Run(":8080")
}