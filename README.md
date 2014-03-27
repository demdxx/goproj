goproj alpha
============

Go language project helper

    @autor Dmitry Ponomarev <demdxx@gmail.com> 2013
    @version 1.0.0alpha
    @licese MIT

## Install

```sh
git clone git://github.com/demdxx/goproj.git

cd goproj

sudo cp goproj /usr/local/bin/
```

## Project structure

 * bin
 * pkg
 * src
   * <project1>/.goproj
   * <project2>/.goproj
   * ...
   * <projectN>/.goproj
 * .gosolution

### Solution file `.gosolution`

```python
{
  '<project name and path same>': {
    'git': '<repository uri>', # Optional
    # OR
    'hg': '<repository uri>', # Optional
  }
}
```

### Project file `.goproj`

```python
# Example revel web project application
{
  "project": "{projectname}",
  "version": "0.0.1",

  # Current project settings
  "run": "{solutionpath}/revel run {app}",
  "build": False,

  # Dependencies
  "deps": {
    "github.com/robfig/revel": {
      "build": "{go} build -o revel {app}/revel", # Run after get
    },
  },
}
```

## Using

All commands are needed packages automatically, so this option can be omitted. Start any commands can be from any subcatalog of project, all settings are set automatically.

### Init new project
```sh
cd project_folder
goproj init <project name or repository URI>
```

### Compile and install packages and dependencies
```sh
goproj install or goproj install <package without src>
```

### Run
```sh
goproj run [target]
```

### Testing
```sh
goproj test [package, ...] [-flags]
```
