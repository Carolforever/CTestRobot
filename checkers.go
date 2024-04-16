package main

import (
	//"io/ioutil"
	"database/sql"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	//"golang.org/x/tools/go/packages"
	//"encoding/csv"
	//"path/filepath"
	"bufio"
	"fmt"
	"log"
	"os"
	//"strconv"
)

// check all the .c file in specified directory
func CheckAll(config Config) {
	StaticAnalysis(config)
	SqlDeal(config)
}

func StaticAnalysis(config Config) {
	run_conf(config)
	CheckSmatch(config)
	CheckCppcheck(config)
}

func run_conf(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir: failed for :", err)
	}

	output, err := exec.Command("python3", robot_dir + "/py_scripts/run_conf.py", robot_dir+"/projects/"+config.Proj_Name, config.Autoconf_Cmd, config.Configure_Cmd).Output()
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

	errStr, err := RunCommand(robot_dir, "cppcheck", robot_dir+"/projects/"+config.Proj_Name)
	if err != nil {
		log.Println("Cpp_Check failed for :", err)
	}
	WriteFile("\ncppcheck:\n", robot_dir+"/result/"+config.Proj_Name+".txt")
	WriteFile(errStr, robot_dir+"/result/"+config.Proj_Name+".txt")
}

func CheckSmatch(config Config) {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	smatch_dir := robot_dir + "/smatch/smatch"
	cgcc_dir := robot_dir + "/smatch/cgcc"
	//create .smatch file for every .c file
	output, err := exec.Command("python3", robot_dir + "/py_scripts/smatch_check.py", robot_dir+"/projects/"+config.Proj_Name, config.Make_Cmd, smatch_dir, cgcc_dir).Output()
	if err != nil {
		log.Println("Smatch_Check failed for :", err)
	}
	result := string(output)
	log.Println(result)

	WriteFile("\nsmatch check :\n", robot_dir+"/result/"+config.Proj_Name+".txt")
	MergeFile(robot_dir+"/projects/"+config.Proj_Name, robot_dir+"/result/"+config.Proj_Name+".txt")
}

func CheckDebian() {
	robot_dir, err := os.Getwd()
	if err != nil {
		log.Println("get current dir failed for :", err)
	}

	db, err := sql.Open("sqlite3", robot_dir + "/result/debian.db")
	if err != nil {
		log.Println("init sqlite failed for :", err)
	}

	rows, err := db.Query("SELECT DISTINCT name, c FROM packages WHERE header > 0 AND c > (\"all\" * 0.4) AND c < 250000")
	if err != nil {
		log.Println("sqlite query failed for :", err)
	}

	var name string
	var c_code int

	for rows.Next() { 
		err = rows.Scan(&name, &c_code)
		if err != nil {
			log.Println("sqlite scan failed for :", err)
		}
	
		cmd := exec.Command("sudo", "apt-cache", "showsrc", name)
		output, _ := cmd.CombinedOutput()
		if strings.Contains(string(output), "Unable to locate package") {
			continue
		}

		pac_basename := name
		pac_name := ""
		
		output, err := exec.Command("python3", robot_dir + "/py_scripts/debian_check.py", robot_dir, pac_basename).Output()
		if err != nil {
			log.Println("Debian check failed for :", err)
		}
		result := string(output)
		log.Println(result)

		cmd = exec.Command("sudo", "chmod", "-R", "777", robot_dir)
		cmd.Run()

		dir, err := os.ReadDir(robot_dir + "/projects")
		if err != nil {
			log.Println("Read_Dir failed for :", err)
		}

		for _, entry := range dir {
			name := entry.Name()
			if entry.IsDir() && !strings.HasSuffix(name, ".git") {
				pac_name = name
				break
			}
		}

		WriteFile("\n" + pac_name + "\tC_code:" + fmt.Sprint(c_code) + "\n", robot_dir + "/result/debian.txt")
		MergeFile(robot_dir + "/projects/" + pac_name, robot_dir + "/result/debian.txt")

		RunCommand(robot_dir + "/projects", "sudo", "rm", "-rf", "*")
	}

	output, err := exec.Command("python3", robot_dir + "/py_scripts/secondary_treatment.py", robot_dir).Output()
	if err != nil {
		log.Println("secondary treatment failed for :", err)
	}
	result := string(output)
	log.Println(result)
}

// merge all the .smatch file into one
func MergeFile(Proj_DIR string, outPath string) {
	dir, err := os.ReadDir(Proj_DIR)
	if err != nil {
		log.Println("Read_Dir failed for :", err)
	}

	for _, entry := range dir {
		name := entry.Name()
		if entry.IsDir() && !strings.Contains(name, ".git") {
			MergeFile(Proj_DIR+"/"+name, outPath)
			continue
		}
		if strings.Contains(name, ".smatch") {
			path := Proj_DIR + "/" + name
			fileData := ReadFile(path)
			WriteFile(fileData, outPath)

			_, err := RunCommand(Proj_DIR, "sudo", "rm", "-rf", path)
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

// 检查表是否存在
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

// 将检测结果存入数据库
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
		) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci`)
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
