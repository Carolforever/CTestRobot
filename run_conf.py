import sys
import subprocess as sp

def run_command(cmd):
    cp = sp.run(cmd, shell=True, capture_output=True, encoding='utf-8')
    if cp.returncode != 0:
        error = f"""Something went wrong with the command: [{cmd}]:{cp.stderr}"""
        raise Exception(error)
    return cp.stdout, cp.stderr

def main(): 
    Proj_DIR = sys.argv[1]
    autoconf_cmd = sys.argv[2]
    configure_cmd = sys.argv[3]
        
    if autoconf_cmd != "":
        cmd = "cd " + Proj_DIR + " && " + autoconf_cmd 
        run_command(cmd)
    
    if configure_cmd != "":
        cmd = "cd " + Proj_DIR + " && " + configure_cmd
        run_command(cmd)

main()