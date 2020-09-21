// +build !darwin

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

const version = "0.0.8"

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

func main() {
	lastNPtr := flag.Int("last", 15, "keep the last specified number of values")
	formatPtr := flag.String("format", "", "regex format for the values.  Accepts NUMBER and UUID as shortcuts.")
	versionPtr := flag.Bool("version", false, "print version number and exit")
	filterPtr := flag.String("filter", "", "regex filter")
	valuePtr := flag.String("value", "{{VALUE}}", "transforms the value before passing it on.")
	ansiPtr := flag.Bool("ansi", false, "support ansi colors in value transformation")
	reversePtr := flag.Bool("reverse", false, "list values in reverse order")

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

		lines = replaceAllPlaceholders(lines, *valuePtr, "VALUE")

		if *ansiPtr {
			lines = handleAnsiCodes(lines)
		}

		listValues(lines)
	default:
		printUsage()
	}

	return
}

func replaceAllPlaceholders(values []string, template string, placeholder string) []string {
	replaced := make([]string, len(values))
	for i, line := range values {
		replaced[i] = replacePlaceholder(template, placeholder, line)
	}
	return replaced
}

func handleAnsiCodes(values []string) []string {
	re := regexp.MustCompile("(?i)\\\\(e|033|x1b)")
	replaced := make([]string, len(values))
	for i, line := range values {
		replaced[i] = re.ReplaceAllString(line, "\x1B")
	}
	return replaced
}

func replacePlaceholder(input string, placeholder string, replacementValue string) string {
	return strings.Replace(input, "{{"+placeholder+"}}", replacementValue, -1)
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
