Goproj 2.x.x
============

[![Gratipay](http://img.shields.io/gratipay/demdxx.svg)](https://gratipay.com/demdxx/)

Golang helps organize and manage projects in the language GO, provides a description and deployment of applications and their dependencies.
NOTE: Small experementall revision for simple projects

    @autor Dmitry Ponomarev <demdxx@gmail.com> 2013 – 2015
    @version 2.1.2
    @licese CC-BY-4,0

Install
=======

```sh
go get -u github.com/demdxx/goproj/cmd/goproj
```

Commands
========

**help** – show help or help [command]

```sh
goproj help
```

## Create new project

```sh
cd {solution-folder}
goproj init {project-name}

# From repository
goproj init git:http://github.com/Undefined/project [project-name]
```

## **.goproj** file

```yaml
name: project
version: 0.0.1
description: Look in the README.md

# Dependencies
deps:
  - github.com/go-martini/martini
  - github.com/martini-contrib/sessions
  - github.com/martini-contrib/oauth2
  - github.com/gopk/config
  - github.com/gopk/templates

# Custom commands
cmd:
  static: "cd '{fullpath}' && gulp build"
  run: "{go} run '{fullpath}/cmd/control/main.go' -basedir '{fullpath}' --debug"
```

## Download and install packages and dependencies

```sh
goproj get [-u] [project name]
```

## Compile packages and dependencies

```sh
goproj build [flags] [project name]
```

## Compile and run Go program

```sh
goproj run [flags] [project name]

# OR run in auto restart mode
# Daemon to observe for a files in your project and automatically restart current command

goproj run [-observer] [project name]
```

TODO
====

 * Add command `goproj init solution`

License
=======

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" property="dct:title">Goproj</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.<br />Based on a work at <a xmlns:dct="http://purl.org/dc/terms/" href="http://github.com/demdxx/goproj" rel="dct:source">http://github.com/demdxx/goproj</a>.
