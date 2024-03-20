# CTestRobot



## Installation



### Install Go

Download the golang package (`go1.20.5.linux-amd64.tar.gz`) and execute the following command:

```
# rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
```



### Install Smatch

```
# apt-get install gcc make sqlite3 libsqlite3-dev libdbd-sqlite3-perl libssl-dev libtry-tiny-perl
# git clone git@github.com:error27/smatch.git
# cd smatch
# make
```



### Install Cppcheck

```
# apt-get install cppcheck
```



## Complete config.json

Fill in the config.json file with sensitive information .



## Building

```
git clone git@github.com:Carolforever/CTestRobot.git
cd CTestRobot
make
./CTestRobot -config config.json
```