package main

import (
)

type Config struct {
	Autoconf_Cmd string `json:"autoconf_cmd"`
	Configure_Cmd string `json:"configure_cmd"`
	Make_Cmd string `json:"make_cmd"`
	Proj_Name string `json:"proj_name"`
	Mysql_Info string `json:"mysql_info"`
}

var Res =  make(chan string)

func main() {
	//botinit()
	//router()
	CheckDebian()
}