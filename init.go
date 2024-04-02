package main

import (
	"log"
	"os"
	"encoding/json"
	"flag"
	"os/exec"
	"bytes"
	//"strings"
	//"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
	//"net/http"
	//"io"
    //"fmt"
)

var (
	flagConfig = flag.String("config", "", "configuration file")
	//flagDebug  = flag.Bool("debug", false, "dump all the logs")
)

func parseConfig(configFile string, user Config) Config {
	var config Config
	bytes, _ := os.ReadFile(configFile)
	value1, _ := sjson.Set(string(bytes), "configure_cmd", user.Configure_Cmd)
	value2, _ := sjson.Set(value1, "make_cmd", user.Make_Cmd)
	value3, _ := sjson.Set(value2, "proj_name", user.Proj_Name)
	value4, _ := sjson.Set(value3, "mysql_info", user.Mysql_Info)
	err := os.WriteFile(configFile, []byte(value4), 0777)
	if err != nil {
		log.Println("Modified config.json failed :", err)
	}
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
	//default make_cmd
	if config.Make_Cmd == "" {
		config.Make_Cmd = "make"
	}
	return config
}

func RunCommand(Dir string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  
	cmd.Stderr = &stderr
	if Dir != "" {
		cmd.Dir = Dir
	}
	err := cmd.Run()
	errStr := stderr.String()
	return errStr, err
}

func botinit() {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	_, err = RunCommand(robot_dir, "ls", "-l", "result")
	if err != nil {
		_, err = RunCommand(robot_dir, "mkdir", "result")
		if err != nil {
			log.Println("mkdir result: failed for :", err)
		}
	}

	_, err= RunCommand(robot_dir, "ls", "-l", "projects")
	if err != nil {
		_, err = RunCommand(robot_dir, "mkdir", "projects")
		if err != nil {
			log.Println("mkdir projects: failed for :", err)
		}
	}

	_, err = RunCommand(robot_dir, "ls", "-l", "smatch")
	if err != nil {
		_, err = RunCommand(robot_dir, "git", "clone", "git@github.com:error27/smatch.git")
		if err != nil {
			log.Println("clone smatch: failed for :", err)
		}

		_, err = RunCommand(robot_dir + "/smatch", "make")
		if err != nil {
			log.Println("smatch make: failed for :", err)
		}
	}
}