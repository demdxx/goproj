Goproj 2.x.x
============

[![Gratipay](http://img.shields.io/gratipay/demdxx.svg)](https://gratipay.com/demdxx/)

Golang helps organize and manage projects in the language GO, provides a description and deployment of applications and their dependencies.
NOTE: Small experementall revision for simple projects

    @autor Dmitry Ponomarev <demdxx@gmail.com> 2013 – 2015
    @version 2.1.1
    @licese CC-BY-4,0

Install
=======

```sh
go get github.com/demdxx/goproj && go install
```

Commands
========

**help** – show help or help [command]

```sh
goproj help
```

Create new project
------------------

```sh
cd {solution-folder}
goproj init {project-name}

# From repository
goproj init git:http://github.com/Undefined/project [project-name]
```

Download and install packages and dependencies
----------------------------------------------

```sh
goproj get
```

Compile packages and dependencies
---------------------------------

```sh
goproj build [flags]
```

Compile and run Go program
--------------------------

```sh
goproj run [project names] [flags]

# OR run in auto restart mode
# Daemon to observe for a files in your project and automatically restart current command

goproj run [project names] -observer
```

License
=======

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" property="dct:title">Goproj</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.<br />Based on a work at <a xmlns:dct="http://purl.org/dc/terms/" href="http://github.com/demdxx/goproj" rel="dct:source">http://github.com/demdxx/goproj</a>.
