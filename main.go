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
}

type File struct{
	Path string `json:"path"`
}


func main() {
	botinit()
	router()
}