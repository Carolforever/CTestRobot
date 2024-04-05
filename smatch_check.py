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
    make_cmd = sys.argv[2]
    smatch_dir = sys.argv[3]
    cgcc_dir = sys.argv[4]
        
    #make CHECK="/home/lsc20011130/CTestRobot/smatch/smatch --full-path --file-output" CC="/home/lsc20011130/CTestRobot/smatch/cgcc"
    cmd = "cd " + Proj_DIR + " && " + make_cmd + " CHECK=\"" + smatch_dir + " --full-path" + " --file-output\" " + "CC=\"" + cgcc_dir + "\"" 
    run_command(cmd)

main()