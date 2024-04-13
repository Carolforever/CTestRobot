# CTestRobot

This project uses smatch and cppcheck to statically analyze C projects that are not part of the Linux Kernel. And I'm currently working on the static analysis of Debian C packages.

## Installation


### Install Go

Wsl2 Debian12(bookworm) / Ubuntu 20.04.3 LTS

Golang version : go1.20.5


### Install dependency

```
# apt-get install gcc make sqlite3 libsqlite3-dev libdbd-sqlite3-perl libssl-dev libtry-tiny-perl dpkg-dev cppcheck Sparse
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
# see the page in localhost:8080/lsc, fill in the information and submit
```



## Input Explanation
```
# autoconf_cmd, configure_cmd and make_cmd， there's nothing to explain, just for the project to build. If the project doesn't require autoconf_cmd and configure_cmd, just don't input
# proj_name will be used as the name of the project's test result file in /CTestRobot/result
# mysql_info, follow the template "username:passwd@tcp(localhost:3306)/database_name"
```