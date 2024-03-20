import sys
import os
import stat

def main(): 
    Proj_DIR = sys.argv[1]
    make_cmd = sys.argv[2]
    smatch_dir = sys.argv[3]
    cgcc_dir = sys.argv[4]
 
    #sudo make CHECK="/home/lsc20011130/smatch/smatch --file-output" CC="/home/lsc20011130/smatch/cgcc"
    cmd = "cd " + Proj_DIR + " && " + make_cmd + " CHECK=\"" + smatch_dir + " --file-output\" " + "CC=\"" + cgcc_dir + "\"" 
    os.system(cmd)

main()