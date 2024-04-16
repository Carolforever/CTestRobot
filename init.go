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

//处理config文件
func parseConfig(configFile string, user Config) Config {
	var config Config
	bytes, _ := os.ReadFile(configFile)
	value1, _ := sjson.Set(string(bytes), "autoconf_cmd", user.Autoconf_Cmd)
	value2, _ := sjson.Set(value1, "configure_cmd", user.Configure_Cmd)
	value3, _ := sjson.Set(value2, "make_cmd", user.Make_Cmd)
	value4, _ := sjson.Set(value3, "proj_name", user.Proj_Name)
	value5, _ := sjson.Set(value4, "mysql_info", user.Mysql_Info)
	err := os.WriteFile(configFile, []byte(value5), 0777)
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
	//default value
	if config.Make_Cmd == "" {
		config.Make_Cmd = "make"
	}
	return config
}

//执行command模板
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

//初始化bot
func botinit() {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	_, err = RunCommand(robot_dir, "ls", "-l", "result")
	if err != nil {
		_, err = RunCommand(robot_dir, "mkdir", "result")
		if err != nil {
			log.Println("mkdir result failed for :", err)
		}
	}

	_, err= RunCommand(robot_dir, "ls", "-l", "projects")
	if err != nil {
		_, err = RunCommand(robot_dir, "mkdir", "projects")
		if err != nil {
			log.Println("mkdir projects failed for :", err)
		}
	}

	_, err = RunCommand(robot_dir, "ls", "-l", "smatch")
	if err == nil {
		_, err = RunCommand(robot_dir, "rm", "-rf", "smatch")
		if err != nil {
			log.Println("rm smatch failed for :", err)
		}
	}
	_, err = RunCommand(robot_dir, "git", "clone", "git@github.com:error27/smatch.git")
	if err != nil {
		log.Println("clone smatch failed for :", err)
	}

	_, err = RunCommand(robot_dir + "/smatch", "make")
	if err != nil {
		log.Println("smatch make failed for :", err)
	}

	_, err = RunCommand(robot_dir, "ls", "-l", "cloc_debian")
	if err == nil {
		_, err = RunCommand(robot_dir, "rm", "-rf", "cloc_debian")
		if err != nil {
			log.Println("rm cloc_dibian failed for :", err)
		}
	}

	_, err = RunCommand(robot_dir, "git", "clone", "git@github.com:hust-open-atom-club/cloc-debian.git")
	if err != nil {
		log.Println("clone cloc_debian failed for :", err)
	}
}