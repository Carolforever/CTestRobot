package main

import (
	//"github.com/gin-gonic/gin"
	//"net/http"
    //"fmt"
)

type Config struct {
	Configure_Cmd string `json:"configure_cmd"`
	Make_Cmd string `json:"make_cmd"`
	Proj_Name string `json:"proj_name"`
	Mysql_Info string `json:"mysql_info"`
	//user:password@tcp(localhost:3306)/database
}

var Res =  make(chan string)

func main() {
	botinit()
	router()
}