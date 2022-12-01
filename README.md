# histkeep

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

history manager keeps a mru history of values in a flat file and helps you manage them

## install

Direct downloads are available through the [releases page](https://github.com/bbeardsley/histkeep-cli/releases/latest).

## usage

```
Usage
    histkeep [options] <command> <file> <command arguments...>
Version
    0.0.9
Options
  -aarg string
    	item arg in alfred. {{VALUE}} is replaced with the item value. (default "{{VALUE}}")
  -acopy string
    	text to copy in alfred when item in filter is copied.  {{VALUE}} is replaced with the item value.
  -agvar value
    	name=value to be passed to alfred as a global variable.  Parameter can be specified multiple times for multiple variables.
  -aicon string
    	filename of icon to show in alfred for each item in filter. {{VALUE}} is replaced with item value.
  -aitem value
    	item to include in alfred list. Parameter can be specified multiple times for multiple items
  -alfred
    	output Alfred JSON list
  -amod value
    	(alt|cmd|ctrl|fn):(var|valid|arg|subtitle|icon):(value|name=value)
  -ansi
    	support ansi colors in value transformation
  -asubtitle string
    	subtitle to display for the item in alfred.  {{VALUE}} is replaced with the item value.
  -atitle string
    	item title in alfred. {{VALUE}} is replaced with the item value. (default "{{VALUE}}")
  -avar value
    	name=value to be passed to alfred.  {{VALUE}} is replaced with item value in value.  Parameter can be specified multiple times for multiple variables.
  -filter string
    	regex filter
  -format string
    	regex format for the values.  Accepts NUMBER and UUID as shortcuts.
  -last int
    	keep the last specified number of values (default 15)
  -reverse
    	list values in reverse order
  -value string
    	transforms the value before passing it on. (default "{{VALUE}}")
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
