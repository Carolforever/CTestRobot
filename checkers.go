package main

import (
	"bytes"
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
func CheckAll(config Config) string {
	result := StaticAnalysis(config)
	return result
}

func StaticAnalysis(config Config) string {
	result := "smatch check:" + "\n"
	smatch_err, checksmatch := CheckSmatch(config)
	if smatch_err {
		result += checksmatch
	}
	
	result = result + "\ncppcheck check:\n\n"
	cppcheck_err, cppcheck := CheckCppcheck(config)
	if cppcheck_err {
		result += cppcheck
	}
	
	return result
}


func CheckCppcheck(config Config) (bool, string) {
	cmd := exec.Command("cppcheck", config.Proj_Dir, "--enable=warning")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  
	cmd.Stderr = &stderr
	err := cmd.Run()
	errStr := stderr.String()
	if err != nil {
		log.Println("Cpp_Check: failed for :", err)
	}
	WriteFile("\ncppcheck:\n\n", "./result/" + config.Proj_Name + ".txt")
	WriteFile(errStr, "./result/" + config.Proj_Name + ".txt")
	return true, errStr
}


func CheckSmatch(config Config) (bool, string) {
	smatch_dir := config.Smatch_Dir + "/smatch"
	cgcc_dir := config.Smatch_Dir + "/cgcc"
	//create .smatch file for every .c file
	output, err := exec.Command("python", "check.py", config.Proj_Dir, config.Make_Cmd, smatch_dir, cgcc_dir).Output()
	if err != nil {
		log.Println("Smatch_Check: failed for :", err)
	}
	result := string(output)
	log.Println(result)

	MergeFile(config.Proj_Dir, "./result/" + config.Proj_Name + ".txt", "smatch check:\n\n")

	dir, err := os.ReadDir(config.Proj_Dir)
	if err != nil {
		log.Println("Dir Read: failed for :", err)
	}
	for _, file := range dir {
		if strings.Contains(file.Name(), ".smatch") {
			cmd := exec.Command("rm", "-rf", file.Name())
			cmd.Dir = config.Proj_Dir
			err := cmd.Run()
			if err != nil {
				log.Println("Dir clean: failed for :", err)
			}
		}
	}

	return true, ReadFile("./result/" + config.Proj_Name + ".txt")
}

func MergeFile(Proj_DIR string, outPath string, HeadLine string) {
	WriteFile(HeadLine, outPath)
	dir, err := os.ReadDir(Proj_DIR)
	if err != nil {
		log.Println("Read_Dir failed for :", err)
	}

	for _, entry := range dir {
		name := entry.Name()
		/*
		if entry.IsDir() {
			MergeFile(name, outPath)
			continue
		}
		*/
		if strings.Contains(name, ".smatch") {
			path := Proj_DIR + "/" + name
			fileData := ReadFile(path)
			WriteFile(fileData, outPath)
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



