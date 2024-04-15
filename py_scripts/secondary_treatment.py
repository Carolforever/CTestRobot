import sys
import os
import subprocess as sp
import sqlite3

def run_command(filename, cmd):
    cp = sp.run(cmd, shell=True, capture_output=True, encoding='utf-8', timeout=480)
    if cp.returncode != 0:
        error = f"[{filename}] : Something went wrong with the command: [{cmd}]"
        print(error)
    return cp.stdout, cp.stderr

def merge(robot_dir, dir):
    files = os.listdir(dir)
    for fi in files:
        fi_d = os.path.join(dir, fi)            
        if os.path.isdir(fi_d):
            merge(robot_dir, fi_d)                  
        else:
            if fi_d.endswith('.smatch'):
                with open(fi_d, 'r') as file:
                    content = file.read()
                    with open(robot_dir + "/result/debian.txt", 'a', encoding="utf-8") as output_file:
                        output_file.write(content)
                        output_file.close()
                file.close()

def remove_duplicate_lines(robot_dir):
    with open(robot_dir + "/result/simplified_debian.txt", 'r') as file:
        lines = file.readlines()
        with open(robot_dir + "/result/only_error_debian.txt", 'w', encoding="utf-8") as output_file:
            pre_line = ""
            for i in range(len(lines)):
                line = lines[i]
                if line.strip() == "" and pre_line.strip() == "":
                    continue
                output_file.write(line)
                pre_line = line
        output_file.close()
    file.close()
    
    cmd = "cd " + robot_dir + "/result && sudo rm -rf error_debian.txt"
    run_command("clean", cmd)
    cmd = "cd " + robot_dir + "/result && sudo rm -rf simplified_debian.txt"
    run_command("clean", cmd)

def improve_result(robot_dir):
    with open(robot_dir + "/result/debian.txt", 'r') as file:
        pre_line = ""
        tested_code = 0
        tested_package = 0
        
        for line in file:
            if ("warn" in line or "error" in line) and "C_code" in pre_line:
                code_index = pre_line.index("C_code:") + len("C_code:")
                code_number = int(pre_line[code_index:].strip())
                
                tested_code += code_number
                tested_package += 1
                
            
            if "C_code" in line or "error" in line or line.strip() == "":
                with open(robot_dir + "/result/error_debian.txt", 'a+', encoding="utf-8") as file:
                    file.write(line)
                    file.close()
            
            pre_line = line
          
        with open(robot_dir + "/result/error_debian.txt", "a+", encoding="utf-8") as f:
            f.write(f"\nTested C code: {tested_code}\n")
            f.write(f"Tested packages: {tested_package}\n")
            f.close()
                
        file.close() 
    
    with open(robot_dir + "/result/error_debian.txt", 'r') as file:
        lines = file.readlines()
        with open(robot_dir + "/result/simplified_debian.txt", 'w', encoding="utf-8") as output_file:
            pre_line = ""
            for i in range(len(lines)):
                line = lines[i]
                if "C_code" in line and pre_line == "\n" and i+1 < len(lines) and lines[i+1] == "\n":
                    output_file.write("\n")
                else:
                    output_file.write(line)
                pre_line = line
        output_file.close()
    file.close()

    remove_duplicate_lines(robot_dir)
    
    
def main(): 
    robot_dir = sys.argv[1]
    #robot_dir = "/home/lsc/CTestRobot"
    
    smatch_dir = robot_dir + "/smatch/smatch --full-path --file-output"
    cgcc_dir = robot_dir + "/smatch/cgcc"
    smatch_check_cmd = "CHECK=\"" + smatch_dir + "\" CC=\"" + cgcc_dir + "\""
    
    temp = sys.stdout
    f = open(robot_dir + "/result/debian_log.txt", 'a+')
    sys.stdout = f
    
    print("\n\n\n[Debian] : Starting secondary treatment...\n\n\n")
    
    with open(robot_dir + "/result/debian_log.txt", 'r') as file:
        for line in file:
            if "sudo apt-get build-dep -y" in line:
                package = line[line.index('[') + 1:line.index(']')]
                cmd = "cd " + robot_dir + "/projects && sudo apt-get source " + package
                run_command(package, cmd)
    
                filename = ""
                dir_list = os.listdir(robot_dir + "/projects")
                for cur_file in dir_list:
                    path = os.path.join(robot_dir + "/projects", cur_file)
                    if os.path.isdir(path):
                        filename = cur_file
                    
                cmd = "cd " + robot_dir + "/projects/" + filename + " && sudo chmod -R 777 debian"
                run_command(filename, cmd)
                
                build_system_c = True
                override_dh_auto_build = False
                dh_auto_build_changed = False
            
                with open(robot_dir + "/projects/" + filename + "/debian/rules", 'r') as file:
                    for line in file:
                        if "--buildsystem=" in line:
                            build_system_c = False
                            print(f"[{filename}] : Build system is not suitable")
                            file.close()
                            break
                    
                if build_system_c == True: 
                    file_data = ""
                    with open(robot_dir + "/projects/" + filename + "/debian/rules", "r", encoding="utf-8") as f:
                        for line in f:
                            if "override_dh_auto_build:" in line:
                                override_dh_auto_build = True
                            
                            elif dh_auto_build_changed == False and "$(MAKE)" in line:
                                line = line.replace("$(MAKE)", "$(MAKE) " + smatch_check_cmd + " ")
                                dh_auto_build_changed = True
                            
                            elif dh_auto_build_changed == False and "dh_auto_build --" in line:
                                line = line.replace("dh_auto_build --", "dh_auto_build -- " + smatch_check_cmd + " ")
                                dh_auto_build_changed = True
                            
                            elif dh_auto_build_changed == False and "dh_auto_build" in line and "override_dh_auto_build:" not in line:
                                line = line.replace("dh_auto_build", "dh_auto_build -- " + smatch_check_cmd + " ")
                                dh_auto_build_changed = True
                            
                            file_data += line
                    
                    with open(robot_dir + "/projects/" + filename + "/debian/rules", "w", encoding="utf-8") as f:
                        f.write(file_data)
                        f.close()
            
                    if override_dh_auto_build == True and dh_auto_build_changed == True:
                        print(f"[{filename}] : dh_auto_build exists and override...")
                        cmd = "cd " + robot_dir + "/projects/" + filename + " && sudo dpkg-buildpackage -us -uc"
                        print(f"[{filename}] : Running dpkg-buildpackage...")
                        run_command(filename, cmd)
                
                    elif override_dh_auto_build == False and dh_auto_build_changed == False:
                        print(f"[{filename}] : dh_auto_build doesn't exist, adding override_dh_auto_build...")
                        with open(robot_dir + "/projects/" + filename + "/debian/rules", "a", encoding="utf-8") as f:
                            f.write("\noverride_dh_auto_build:\n")
                            f.write("\tdh_auto_build -- " + smatch_check_cmd + "\n")
                            f.close()
                        
                        cmd = "cd " + robot_dir + "/projects/" + filename + " && sudo dpkg-buildpackage -us -uc"
                        print(f"[{filename}] : Running dpkg-buildpackage...")
                        run_command(filename, cmd)
                
                    elif override_dh_auto_build == True and dh_auto_build_changed == False:
                        print(f"[{filename}] : Can't judge, skipping smatch check...")      
                else:               
                    print(f"[{filename}] : Build system is not suitable, skipping smatch check...")                   
            
                with open(robot_dir + "/result/debian.txt", 'a', encoding="utf-8") as file:
                    conn = sqlite3.connect(robot_dir + "/result/debian.db")
                    cursor = conn.cursor()

                    cursor.execute("SELECT DISTINCT c FROM packages WHERE name = ?", (package,))
                    result = cursor.fetchone()

                    conn.close()
                    
                    c_code = result[0]
                    file.write("\n" + filename + "\tC_code:" + str(c_code) + "\n")
                    
                merge(robot_dir, robot_dir + "/projects/" + filename)
            
                cmd = "cd " + robot_dir + "/projects && sudo rm -rf *"
                run_command("fflush", cmd)
            
        file.close()
    
    f.close()
    sys.stdout = temp
    
    improve_result(robot_dir)
        
main()