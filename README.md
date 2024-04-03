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



### Building

```
git clone git@github.com:Carolforever/CTestRobot.git
cd CTestRobot
make
./CTestRobot -config config.json
```



### Use

```
# Put the projects to be tested into /CTestRobot/projects and finish all the required pre-processing steps before ./configure(e.g. installing dependencies required by the project)
# localhost:8080/lsc
```



## Input Explanation
```
# configure_cmd and make_cmd， there's nothing to explain, just for the project to build.
# proj_name will be used as the name of the project's test result file in /CTestRobot/result
# mysql_info, follow the template "username:passwd@tcp(localhost:3306)/database_name"
```