package main

import (
	//"io/ioutil"
	"os/exec"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	SqlDeal(config)
}


func StaticAnalysis(config Config)  {
	run_conf(config)
	CheckSmatch(config)
	CheckCppcheck(config)
}

func run_conf(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	output, err := exec.Command("python", "run_conf.py", robot_dir + "/projects/" + config.Proj_Name, config.Autoconf_Cmd, config.Configure_Cmd).Output()
	if err != nil {
		log.Println("run_conf failed for :", err)
	}
	result := string(output)
	log.Println(result)
}

func CheckCppcheck(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	errStr, err := RunCommand(robot_dir, "cppcheck", robot_dir + "/projects/" + config.Proj_Name)
	if err != nil {
		log.Println("Cpp_Check failed for :", err)
	}
	WriteFile("\ncppcheck:\n", robot_dir + "/result/" + config.Proj_Name + ".txt")
	WriteFile(errStr, robot_dir + "/result/" + config.Proj_Name + ".txt")
}


func CheckSmatch(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	smatch_dir := robot_dir + "/smatch/smatch"
	cgcc_dir := robot_dir + "/smatch/cgcc"
	//create .smatch file for every .c file
	output, err := exec.Command("python", "smatch_check.py", robot_dir + "/projects/" + config.Proj_Name, config.Make_Cmd, smatch_dir, cgcc_dir).Output()
	if err != nil {
		log.Println("Smatch_Check failed for :", err)
	}
	result := string(output)
	log.Println(result)

	WriteFile("\nsmatch check :\n", robot_dir + "/result/" + config.Proj_Name + ".txt")
	MergeFile(robot_dir + "/projects/" + config.Proj_Name, robot_dir + "/result/" + config.Proj_Name + ".txt")
}

//merge all the .smatch file into one
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

			_, err := RunCommand(Proj_DIR, "rm", "-rf", path)
			if err != nil {
				log.Println("Dir clean failed for :", err)
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

//检查表是否存在
func checkTableExists(db *sql.DB, tableName string) (bool, error) {
    rows, err := db.Query("SHOW TABLES")
    if err != nil {
        return false, err
    }
    defer rows.Close()

    for rows.Next() {
        var table string
        if err := rows.Scan(&table); err != nil {
            return false, err
        }
        if table == tableName {
            return true, nil
        }
    }
    return false, nil
}

//将检测结果存入数据库
func SqlDeal(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	db, err := sql.Open("mysql", config.Mysql_Info)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
		
	//新建表存储检测结果
	exists, err := checkTableExists(db, "files")
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = db.Exec(`CREATE TABLE files (
			id int NOT NULL AUTO_INCREMENT, filename varchar(255) NOT NULL, content longtext NOT NULL, created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci`)
		if err != nil {	
			panic(err.Error())
		}
	}

	// 读取文件内容
	filePath := robot_dir + "/result/" + config.Proj_Name + ".txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	Res <- string(content)
		
	// 插入数据到表中
	insertStmt, err := db.Prepare("INSERT INTO files (filename, content) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertStmt.Close()
		
	_, err = insertStmt.Exec(config.Proj_Name, string(content))
	if err != nil {
		panic(err.Error())
	}
}