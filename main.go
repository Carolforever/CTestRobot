package main

import (
	"log"
	"os"
	"encoding/json"
	"flag"
	//"io"
    //"fmt"

)

type Config struct {
	Make_Cmd string `json:"make_cmd"`
	Smatch_Dir string `json:"smatch_dir"`
	Proj_Dir string `json:"proj_dir"`
	Proj_Name string `json:"proj_name"`
}

var (
	flagConfig = flag.String("config", "", "configuration file")
	//flagDebug  = flag.Bool("debug", false, "dump all the logs")
)

func parseConfig(configFile string) Config {
	var config Config
	// open config file
	configFd, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer configFd.Close()
	// parse json file
	dec := json.NewDecoder(configFd)
	// disallow any unknown fields
	dec.DisallowUnknownFields()
	if err = dec.Decode(&config); err != nil {
		log.Fatal(err)
	}
	if config.Make_Cmd == "" {
		config.Make_Cmd = "make"
	}
	if config.Smatch_Dir == "" {
		config.Smatch_Dir = "./smatch"
	}
	return config
}

func main() {
	flag.Parse()
	if *flagConfig == "" {
		log.Fatalf("No config file specified")
	}
	config := parseConfig(*flagConfig)
	result := CheckAll(config)
	log.Println(result)
}