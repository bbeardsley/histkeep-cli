package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bbeardsley/histkeep"
)

const version = "0.0.6"

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage")
	fmt.Fprintln(os.Stderr, "    histkeep [options] <command> <file> <command arguments...>")
	fmt.Fprintln(os.Stderr, "Version")
	fmt.Fprintln(os.Stderr, "    "+version)
	fmt.Fprintln(os.Stderr, "Options")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Commands")
	fmt.Fprintln(os.Stderr, "  help    -> show this help")
	fmt.Fprintln(os.Stderr, "  version -> print version number and exit")
	fmt.Fprintln(os.Stderr, "  add     -> add value")
	fmt.Fprintln(os.Stderr, "  clear   -> clear all values")
	fmt.Fprintln(os.Stderr, "  list    -> list values")
	fmt.Fprintln(os.Stderr, "  remove  -> remove value")
	os.Exit(1)
}

var alfredGlobalVarFlags arrayFlags
var alfredItemVarFlags arrayFlags
var alfredCannedItemFlags arrayFlags

func main() {
	lastNPtr := flag.Int("last", 15, "keep the last specified number of values")
	formatPtr := flag.String("format", "", "regex format for the values.  Accepts NUMBER and UUID as shortcuts.")
	versionPtr := flag.Bool("version", false, "print version number and exit")
	alfredPtr := flag.Bool("alfred", false, "output Alfred JSON list")
	reversePtr := flag.Bool("reverse", false, "list values in reverse order")
	acopyPtr := flag.String("acopy", "", "text to copy in alfred when item in filter is copied.  {{VALUE}} is replaced with the item value.")
	aiconPtr := flag.String("aicon", "", "filename of icon to show in alfred for each item in filter. {{VALUE}} is replaced with item value.")
	asubtitlePtr := flag.String("asubtitle", "", "subtitle to display for the item in alfred.  {{VALUE}} is replaced with the item value.")
	atitlePtr := flag.String("atitle", "{{VALUE}}", "item title in alfred. {{VALUE}} is replaced with the item value.")
	aargPtr := flag.String("aarg", "{{VALUE}}", "item arg in alfred. {{VALUE}} is replaced with the item value.")
	filterPtr := flag.String("filter", "", "regex filter")
	flag.Var(&alfredItemVarFlags, "avar", "name=value to be passed to alfred.  {{VALUE}} is replaced with item value in value.  Parameter can be specified multiple times for multiple variables.")
	flag.Var(&alfredCannedItemFlags, "aitem", "item to include in alfred list. Parameter can be specified multiple times for multiple items")
	flag.Var(&alfredGlobalVarFlags, "agvar", "name=value to be passed to alfred as a global variable.  Parameter can be specified multiple times for multiple variables.")

	flag.Parse()

	if *versionPtr {
		fmt.Println(version)
		os.Exit(0)
	}

	command := strings.TrimSpace(flag.Arg(0))
	file := strings.TrimSpace(flag.Arg(1))
	value := strings.TrimSpace(flag.Arg(2))
	format := buildFormat(*formatPtr)
	hist := histkeep.NewHistKeep(file, *lastNPtr, format)

	switch command {
	case "", "h", "-h", "--h", "/h", "/?", "help", "-help", "--help", "/help":
		printUsage()
	case "version", "-version", "--version", "/version":
		fmt.Println(version)
	case "add":
		if file == "" || value == "" {
			printUsage()
		}

		err := hist.AddValue(value)
		if err != nil {
			log.Fatal(err)
		}
	case "clear":
		if file == "" {
			printUsage()
		}

		err := hist.ClearValues()
		if err != nil {
			log.Fatal(err)
		}
	case "remove":
		if file == "" || value == "" {
			printUsage()
		}

		err := hist.RemoveValue(value)
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		if file == "" {
			printUsage()
		}

		filterFunc := buildFilterFunc(*filterPtr)
		lines, err := hist.GetFilteredValues(filterFunc)
		if err != nil {
			log.Fatal(err)
		}

		if *reversePtr {
			lines = hist.ReverseValues(lines)
		}

		if *alfredPtr {
			a := alfred{
				itemTitle:    *atitlePtr,
				itemArg:      *aargPtr,
				itemSubtitle: *asubtitlePtr,
				iconFilename: *aiconPtr,
				copyText:     *acopyPtr,
				itemVars:     alfredItemVarFlags,
				cannedItems:  alfredCannedItemFlags,
				filter:       *filterPtr,
				format:       format,
				globalVars:   alfredGlobalVarFlags,
			}
			a.list(lines)
		} else {
			listValues(lines)
		}
	default:
		printUsage()
	}

	return
}

func buildFilterFunc(filter string) func(string) bool {
	var filterFunc func(string) bool
	if filter != "" {
		filterRegex, err := regexp.Compile("^(?i)" + filter)
		if err != nil {
			log.Fatal(err)
		}
		filterFunc = func(line string) bool {
			return filterRegex.MatchString(line)
		}
	} else {
		filterFunc = func(line string) bool {
			return true
		}
	}
	return filterFunc
}

func listValues(values []string) {
	for i, line := range values {
		if i != 0 {
			fmt.Println()
		}
		fmt.Print(line)
	}
}

func buildFormat(formatStr string) *regexp.Regexp {
	format := processedNamedFormats(formatStr)
	return regexp.MustCompile("^" + format + "$")
}

func processedNamedFormats(formatStr string) string {
	switch formatStr {
	case "NUMBER":
		return "\\d+"
	case "UUID":
		return "([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}"
	case "":
		return ".*"
	default:
		return formatStr
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
