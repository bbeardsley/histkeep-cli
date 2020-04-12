# histkeep

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

history manager keeps of history of values in a flat file and helps you manage them

## install

Direct downloads are available through the [releases page](https://github.com/bbeardsley/histkeep/releases/latest).

If you have Go installed on your computer just run `go get`.

    go get github.com/bbeardsley/histkeep

## usage

```
Usage
    histkeep [options] <command> <file> <command arguments...>
Version
    0.0.1
Options
  -last int
        keep the last specified number of values (default 15)
  -version
        print version number and exit
Commands
  help    -> show this help
  version -> print version number and exit
  add     -> add value
  clear   -> clear all values
  list    -> list values
  remove  -> remove value
```
