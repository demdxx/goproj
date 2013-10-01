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

## Project file *.goproj*

```python
# Example revel web project application
{
  "project": "{projectname}",
  "version": "0.0.1",

  "deps": {
    "github.com/robfig/revel": {
      "build": "{go} build -o revel {app}/revel", # Run after get
    },
  },

  "apps": [
    # default all from src
    "{projectname}/app"
  ],

  # @TODO: HOOKS
  # "hooks": {"command": {"before": ["script/path", ...]}, "after": ["script/path", ...]},

  # "build": ["arg1", "arg2", ...]
}
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

* Testing
```sh
goproj test [package, ...] [-flags]
```
