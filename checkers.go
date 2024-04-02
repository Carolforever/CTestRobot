package main

import (
	//"io/ioutil"
	"os/exec"
	"strings"
	//"path/filepath"
	"fmt"
	//"crypto/sha256"
	//"strings"
	"bufio"
	"log"
	"os"
	//"strconv"
)

//check all the .c file in specified directory
func CheckAll(config Config) {
	StaticAnalysis(config)
}

func StaticAnalysis(config Config)  {
	CheckSmatch(config)
	CheckCppcheck(config)
}


func CheckCppcheck(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	err, errStr := RunCommand(robot_dir, "cppcheck", robot_dir + "/projects/" + config.Proj_Name, "--enable=warning")
	if err != nil {
		log.Println("Cpp_Check: failed for :", err)
	}
	WriteFile("\ncppcheck:\n", robot_dir + "/result/" + config.Proj_Name + ".txt")
	WriteFile(errStr, robot_dir + "/result/" + config.Proj_Name + ".txt")
}


func CheckSmatch(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	smatch_dir := robot_dir + "/smatch/smatch"
	cgcc_dir := robot_dir + "/smatch/cgcc"
	//create .smatch file for every .c file
	output, err := exec.Command("python", "check.py", robot_dir + "/projects/" + config.Proj_Name, config.Configure_Cmd, config.Make_Cmd, smatch_dir, cgcc_dir).Output()
	if err != nil {
		log.Println("Smatch_Check: failed for :", err)
	}
	result := string(output)
	log.Println(result)

	WriteFile("smatch check :\n", robot_dir + "/result/" + config.Proj_Name + ".txt")
	MergeFile(robot_dir + "/projects/" + config.Proj_Name, robot_dir + "/result/" + config.Proj_Name + ".txt")
}

func MergeFile(Proj_DIR string, outPath string) {
	dir, err := os.ReadDir(Proj_DIR)
	if err != nil {
		log.Println("Read_Dir failed for :", err)
	}

	for _, entry := range dir {
		name := entry.Name()
		if entry.IsDir() && !strings.Contains(name, ".git"){
			MergeFile(Proj_DIR + "/" + name, outPath)
			continue
		}
		if strings.Contains(name, ".smatch") {
			path := Proj_DIR + "/" + name
			fileData := ReadFile(path)
			WriteFile(fileData, outPath)

			err, _ := RunCommand(Proj_DIR, "rm", "-rf", path)
			if err != nil {
				log.Println("Dir clean: failed for :", err)
			}
		}
	}
}

func WriteFile(s string, outPath string) {
	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bufio.NewWriter(file)
	buf.WriteString(s)
	buf.Flush()
}

func ReadFile(readPath string) string {
	by, err := os.ReadFile(readPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(by)
}



