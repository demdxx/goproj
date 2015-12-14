#!/usr/bin/env bash

# Check project environment
goproj env

# Check original go environment
goproj go env

# Load project from github and put to "src/gopm"
goproj init --autoconfig https://github.com/GPMGo/gopm#c855bb10ada00278fbcc494467e573faf2bd9f70 gopm

# Set options for target=test config
goproj options set production --env-os=linux --env-arch=arm

# Add global dependency or dep
goproj dependency add github.com/demdxx/gocast

# Show list of dependencies
goproj dependency list

# Run test for project --verbose
goproj test -v

# Build project gopm
goproj build --taget=production gopm