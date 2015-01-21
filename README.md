Goproj 2.x.x
============

<a href="https://gratipay.com/demdxx/"><img src="//img.shields.io/gratipay/demdxx.svg"></a>

Golang helps organize and manage projects in the language GO, provides a description and deployment of applications and their dependencies.
NOTE: Small experementall revision for simple projects

    @autor Dmitry Ponomarev <demdxx@gmail.com> 2013
    @version 1.0.0alpha
    @licese CC-BY-4,0

Install
=======

```sh
go get github.com/demdxx/goproj && go install
```

Commands
========

```sh
cd {any-dir-in-solution}
```

    flags: # eqvals env vars GO{flag}
      os = mac | linux | freebsd | etc
      gc = off | {number}
      arch = arm | i386 | x86_64 | etc
      -arm -> shortcut for -arch arm

**help** – show help or help [command]

```sh
goproj help
```

**get** – download and install packages and dependencies.

```sh
goproj get
```

**build** – compile packages and dependencies.

```sh
goproj build [flags]
```

**run** – compile and run Go program.

```sh
goproj run [project names] [flags]
```

TODO
====

 * Postanalize dependencies if custom install
 * Add support .godeps
 * Add update environment flags

License
=======

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" property="dct:title">Goproj</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.<br />Based on a work at <a xmlns:dct="http://purl.org/dc/terms/" href="http://github.com/demdxx/goproj" rel="dct:source">http://github.com/demdxx/goproj</a>.
