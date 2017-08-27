#!/bin/sh

## DO NOT RUN SCRIPT DIRECTLY
## Use `make deps-update`

set -e

ifupdated() {
  local lastfile=".git/${2}"
  if [ ! -f "$lastfile" ]; then
    echo "no dep commit file for ${1}"
    make $2
    return
  fi
  local last="$(cat $lastfile)"
  local current="$(git log -n 1 --pretty=format:%h -- ${1})"
  if [ "$current" != "$last" ]; then
    make $2
  fi
}

#ifupdated glide.yaml deps-go
#ifupdated ui/package.json deps-js
#ifupdated ui/semantic deps-css
