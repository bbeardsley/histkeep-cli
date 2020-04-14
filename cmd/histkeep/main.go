package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/bbeardsley/histkeep"
)

const version = "0.0.4"

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

var alfredVarFlags arrayFlags

func main() {
	lastNPtr := flag.Int("last", 15, "keep the last specified number of values")
	formatPtr := flag.String("format", "", "regex format for the values.  Accepts NUMBER and UUID as shortcuts.")
	versionPtr := flag.Bool("version", false, "print version number and exit")
	alfredPtr := flag.Bool("alfred", false, "output Alfred JSON list")
	reversePtr := flag.Bool("reverse", false, "list values in reverse order")
	acopyPtr := flag.String("acopy", "", "text to copy in alfred when item in filter is copied.  {{VALUE}} is replaced with the item value.")
	aiconPtr := flag.String("aicon", "", "filename of icon to show in alfred for each item in filter. {{VALUE}} is replaced with item value.")
	flag.Var(&alfredVarFlags, "avar", "name=value to be passed to alfred.  {{VALUE}} is replaced with item value in both name and value.  Parameter can be specified multiple times for multiple variables.")

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

		lines, err := hist.GetValues()
		if err != nil {
			log.Fatal(err)
		}

		if *reversePtr {
			lines = reverseValues(lines)
		}

		if *alfredPtr {
			listAlfred(lines, *aiconPtr, *acopyPtr, alfredVarFlags)
		} else {
			listValues(lines)
		}
	default:
		printUsage()
	}

	return
}

func listAlfred(values []string, iconFilename string, copyText string, alfredVars arrayFlags) {
	fmt.Println("{\"items\": [")
	for i, line := range values {
		if i != 0 {
			fmt.Print(",")
		}

		fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\"", line, line)
		if iconFilename != "" {
			fmt.Print(", \"icon\": { \"path\": \"")
			fmt.Print(strings.Replace(iconFilename, "{{VALUE}}", line, -1))
			fmt.Print("\"} ")
		}
		if copyText != "" {
			fmt.Print(", \"text\": { \"copy\": \"")
			fmt.Print(strings.Replace(copyText, "{{VALUE}}", line, -1))
			fmt.Print("\"} ")
		}
		if len(alfredVars) > 0 {
			fmt.Print(", \"variables\": {")
			for i, avar := range alfredVars {
				if i != 0 {
					fmt.Print(",")
				}
				parts := strings.Split(avar, "=")
				if len(parts) == 2 {
					fmt.Printf("\"%v\": \"%v\"", strings.Replace(parts[0], "{{VALUE}}", line, -1), strings.Replace(parts[1], "{{VALUE}}", line, -1))
				}
			}
			fmt.Println("}")
		}
		fmt.Println("}")
	}
	fmt.Println("]}")
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
	regex, _ := regexp.Compile("^" + format + "$")
	return regex
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
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type stringValue struct {
	value string
	index int
}

type byIndex []stringValue

func (b byIndex) Len() int           { return len(b) }
func (b byIndex) Less(i, j int) bool { return b[i].index < b[j].index }
func (b byIndex) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func reverseValues(values []string) []string {
	stringValues := make([]stringValue, 0)
	for i := 0; i < len(values); i++ {
		stringValues = append(stringValues, stringValue{values[i], i})
	}
	sort.Sort(sort.Reverse(byIndex(stringValues)))
	items := make([]string, 0)
	for i := 0; i < len(stringValues); i++ {
		items = append(items, stringValues[i].value)
	}
	return items
}
