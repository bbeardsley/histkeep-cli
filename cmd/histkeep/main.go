package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bbeardsley/histkeep"
)

const version = "0.0.2"

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
}

func main() {
	lastNPtr := flag.Int("last", 15, "keep the last specified number of values")
	versionPtr := flag.Bool("version", false, "print version number and exit")

	flag.Parse()

	if *versionPtr {
		fmt.Println(version)
		os.Exit(0)
	}

	command := flag.Arg(0)
	file := flag.Arg(1)
	value := flag.Arg(2)

	hist := histkeep.NewHistKeep(file, *lastNPtr)

	switch command {
	case "", "h", "-h", "--h", "/h", "/?", "help", "-help", "--help", "/help":
		printUsage()
		os.Exit(1)
	case "version", "-version", "--version", "/version":
		fmt.Println(version)
	case "add":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		if value == "" {
			printUsage()
			os.Exit(1)
		}
		err := hist.AddValue(value)
		if err != nil {
			log.Fatal(err)
		}
	case "clear":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		err := hist.ClearValues()
		if err != nil {
			log.Fatal(err)
		}
	case "remove":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		if value == "" {
			printUsage()
			os.Exit(1)
		}
		err := hist.RemoveValue(value)
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		err := hist.ListValues()
		if err != nil {
			log.Fatal(err)
		}
	default:
		printUsage()
	}

	return
}
