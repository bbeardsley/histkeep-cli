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
    0.0.5
Options
  -aarg string
    	item arg in alfred. {{VALUE}} is replaced with the item value. (default "{{VALUE}}")
  -acopy string
    	text to copy in alfred when item in filter is copied.  {{VALUE}} is replaced with the item value.
  -aicon string
    	filename of icon to show in alfred for each item in filter. {{VALUE}} is replaced with item value.
  -aitem value
    	item to include in alfred list. Parameter can be specified multiple times for multiple items
  -alfred
    	output Alfred JSON list
  -asubtitle string
    	subtitle to display for the item in alfred.  {{VALUE}} is replaced with the item value.
  -atitle string
    	item title in alfred. {{VALUE}} is replaced with the item value. (default "{{VALUE}}")
  -avar value
    	name=value to be passed to alfred.  {{VALUE}} is replaced with item value in both name and value.  Parameter can be specified multiple times for multiple variables.
  -filter string
    	regex filter
  -format string
    	regex format for the values.  Accepts NUMBER and UUID as shortcuts.
  -last int
    	keep the last specified number of values (default 15)
  -reverse
    	list values in reverse order
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
