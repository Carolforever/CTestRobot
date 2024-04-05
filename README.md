# CTestRobot



## Installation



### Install Go

Download the golang package (`go1.20.5.linux-amd64.tar.gz`) and execute the following command:

```
# rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
```



### Install Smatch's dependency

```
# apt-get install gcc make sqlite3 libsqlite3-dev libdbd-sqlite3-perl libssl-dev libtry-tiny-perl
```

### Install Cppcheck

```
# apt-get install cppcheck
```



## Building

```
git clone git@github.com:Carolforever/CTestRobot.git
cd CTestRobot
make
```



## Use

```
# Put the projects to be tested into /CTestRobot/projects and finish all the required pre-processing steps before autoconf(e.g. installing dependencies required by the project)
# Makefile of the project to be tested need to use CC == gcc
# start mysql service
# run ./CTestRobot -config config.json in CTestRobot
# see the page in localhost:8080/lsc, and fill in the information
```



## Input Explanation
```
# configure_cmd and make_cmd， there's nothing to explain, just for the project to build. If the project doesn't require configure_cmd, just don't enter it
# proj_name will be used as the name of the project's test result file in /CTestRobot/result
# mysql_info, follow the template "username:passwd@tcp(localhost:3306)/database_name"
```