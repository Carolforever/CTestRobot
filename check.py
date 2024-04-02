import sys
import os
import stat

def main(): 
    Proj_DIR = sys.argv[1]
    configure_cmd = sys.argv[2]
    make_cmd = sys.argv[3]
    smatch_dir = sys.argv[4]
    cgcc_dir = sys.argv[5]
 
    #./configure
    if configure_cmd != "":
        cmd = "cd " + Proj_DIR + " && " + configure_cmd
        os.system(cmd)
        
    #make CHECK="/smatch/smatch --full-path --file-output" CC="/home/lsc20011130/smatch/cgcc"
    cmd = "cd " + Proj_DIR + " && " + make_cmd + " CHECK=\"" + smatch_dir + " --full-path" + " --file-output\" " + "CC=\"" + cgcc_dir + "\"" 
    os.system(cmd)

main()