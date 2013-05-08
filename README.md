goproj alpha
============

Go language project helper

    @autor Dmitry Ponomarev <demdxx@gmail.com> 2013
    @licese MIT

## Install

```sh
git clone git://github.com/demdxx/goproj.git

cd goproj

sudo cp goproj /usr/local/bin/
```
## Using

All commands are needed packages automatically, so this option can be omitted. Start any commands can be from any podkotalaga project, all settings are set automatically.

* Init new project
```sh
cd project_folder
goproj init <project name>
```

* Compile and install packages and dependencies
```sh
goproj install or goproj install <package without src>
```
