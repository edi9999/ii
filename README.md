# godu
[![Build Status](https://travis-ci.org/edi9999/ii.svg?branch=master)](https://travis-ci.org/edi9999/ii)
[![codecov](https://codecov.io/gh/edi9999/ii/branch/master/graph/badge.svg)](https://codecov.io/gh/edi9999/ii)
[![Go Report Card](https://goreportcard.com/badge/github.com/edi9999/ii)](https://goreportcard.com/report/github.com/edi9999/ii)
[![Gitter chat](https://badges.gitter.im/viktomas-godu.png)](https://gitter.im/viktomas-godu)

Find the files that are taking up your space.

<img src="https://media.giphy.com/media/AhMAsxHCOM1Ve/giphy.gif" width="100%" />

Tired of looking like a noob with [Disk Inventory X](http://www.derlien.com/), [Daisy Disk](https://daisydiskapp.com/) or SpaceMonger? Do you want something that
* can do the job
* scans your drive blazingly fast
* works in terminal
* makes you look cool
* is written in Golang
* you can contribute to

??

Well then **look no more** and try out the godu.

## Installation
```
go get -u github.com/edi9999/ii
```

## Configuration
You can specify names of ignored folders in `.goduignore` in your home directory:
```
> cat ~/.goduignore
node_modules
>
```
I found that I could reduce time it took to crawl through the whole drive to 25% when I started ignoring all `node_modules` which cumulatively contain gigabytes of small text files.

The `.goduignore` is currently only supporting whole folder names. PR that will make it work like `.gitignore` is welcomed.

## Usage
```
godu ~
godu -l 100 / # walks the whole root but shows only files larger than 100MB
# godu ~ | xargs rm # use with caution! Will delete all marked files!
```

The currently selected file / folder can be un/marked with the space-key. Upon exiting, godu prinsts all marked files & folders to stdout so they can be further processed (e.g. via the `xargs` command).

Mind you `-l  <size_limit_in_mb>` option is not speeding up the walking process, it just allows you to filter small files you are not interested in from the output. **The default limit is 10MB**.

Use arrows to move around, space to select a file / folder, ESC or CTRL+C to quit
