import subprocess as sp
import os
import sys
import sqlite3

def run_command(filename, cmd):
    cp = sp.run(cmd, shell=True, capture_output=True, encoding='utf-8', timeout=600)
    if cp.returncode != 0:
        error = f"[{filename}] : Something went wrong with the command: [{cmd}]"
        print(error)
        raise Exception(error)
    return cp.stdout, cp.stderr

def search_files(directory):
    c_file_found = False
    h_file_found = False
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".c"):
                c_file_found = True
            elif file.endswith(".h"):
                h_file_found = True
            if c_file_found and h_file_found:
                return True
    return False

def main(): 
    
    robot_dir = sys.argv[1]
    package = sys.argv[2]
    
    smatch_dir = robot_dir + "/smatch/smatch --full-path --file-output"
    cgcc_dir = robot_dir + "/smatch/cgcc"
    smatch_check_cmd = "CHECK=\"" + smatch_dir + "\" CC=\"" + cgcc_dir + "\""
    
    temp = sys.stdout
    f = open(robot_dir + "/result/debian_log.txt", 'a+')
    sys.stdout = f
    
    cmd = "cd ./projects && sudo apt-get source " + package
    run_command(package, cmd)
    
    filename = ""
    dir_list = os.listdir(robot_dir + "/projects")
    for cur_file in dir_list:
        path = os.path.join(robot_dir + "/projects", cur_file)
        if os.path.isdir(path):
            filename = cur_file

    if search_files(os.path.join(robot_dir, "projects", filename)):         
        cmd = "cd ./projects/" + filename + " && sudo chmod -R 777 debian"
        run_command(filename, cmd)
                
        build_system_c = True
        override_dh_auto_build = False
        dh_auto_build_changed = False
            
        with open("./projects/" + filename + "/debian/rules", 'r') as file:
            for line in file:
                if "--buildsystem=" in line:
                    build_system_c = False
                    print(f"[{filename}] : Build system is not suitable")
                    file.close()
                    break
                    
        if build_system_c == True: 
            file_data = ""
            with open("./projects/" + filename + "/debian/rules", "r", encoding="utf-8") as f:
                for line in f:
                    if "override_dh_auto_build:" in line:
                        override_dh_auto_build = True
                            
                    elif "dh_auto_build --" in line:
                        line = line.replace("dh_auto_build --", "dh_auto_build -- " + smatch_check_cmd)
                        dh_auto_build_changed = True
                            
                    elif "dh_auto_build" in line and "override_dh_auto_build:" not in line:
                        line = line.replace("dh_auto_build", "dh_auto_build -- " + smatch_check_cmd)
                        dh_auto_build_changed = True
                            
                    file_data += line
                    
            with open("./projects/" + filename + "/debian/rules", "w", encoding="utf-8") as f:
                f.write(file_data)
                f.close()
            
            if override_dh_auto_build == True and dh_auto_build_changed == True:
                cmd = "sudo apt-get build-dep -y " + package
                run_command(package, cmd)
                
                print(f"[{filename}] : dh_auto_build exists and override...")
                cmd = "cd ./projects/" + filename + " && sudo dpkg-buildpackage -us -uc"
                print(f"[{filename}] : Running dpkg-buildpackage...")
                run_command(filename, cmd)
                
            elif override_dh_auto_build == False and dh_auto_build_changed == False:
                print(f"[{filename}] : dh_auto_build doesn't exist, adding override_dh_auto_build...")
                with open("./projects/" + filename + "/debian/rules", "a", encoding="utf-8") as f:
                    f.write("\noverride_dh_auto_build:\n")
                    f.write("\tdh_auto_build -- " + smatch_check_cmd + "\n")
                    f.close()
                 
                cmd = "sudo apt-get build-dep -y " + package
                run_command(package, cmd) 
                        
                cmd = "cd ./projects/" + filename + " && sudo dpkg-buildpackage -us -uc"
                print(f"[{filename}] : Running dpkg-buildpackage...")
                run_command(filename, cmd)
                
            elif override_dh_auto_build == True and dh_auto_build_changed == False:
                print(f"[{filename}] : Can't judge, skipping smatch check...")      
        
        else:
            print(f"[{filename}] : Build system is not suitable, skipping smatch check...")  
    else:
        print(f"[{filename}] : No .c or .h files found, skipping smatch check...")
    
main()