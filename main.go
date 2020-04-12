package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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

func readLines(path string, ignoreValue string) ([]string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != ignoreValue && line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func limitSlice(lines []string, lastN int) ([]string, error) {
	linesLen := len(lines)
	if linesLen > lastN {
		return lines[linesLen-lastN : linesLen], nil
	}
	return lines, nil
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for i, line := range lines {
		if i != len(lines)-1 {
			fmt.Fprintln(w, line)
		} else {
			fmt.Fprint(w, line)
		}
	}
	return w.Flush()
}

func addValue(file string, value string, lastN int) error {
	lines, err := readLines(file, value)
	if err != nil {
		return err
	}

	lines = append(lines, value)

	lines, err = limitSlice(lines, lastN)
	if err != nil {
		return err
	}

	err = writeLines(file, lines)
	if err != nil {
		return err
	}
	return nil
}

func removeValue(file string, value string) error {
	lines, err := readLines(file, value)
	if err != nil {
		return err
	}

	err = writeLines(file, lines)
	if err != nil {
		return err
	}
	return nil
}

func clearValues(file string) error {
	lines := make([]string, 0)

	err := writeLines(file, lines)
	if err != nil {
		return err
	}

	return nil
}

func listValues(file string, lastN int) error {
	lines, err := readLines(file, "")
	if err != nil {
		return err
	}

	lines, err = limitSlice(lines, lastN)
	if err != nil {
		return err
	}

	for i, line := range lines {
		if i != len(lines)-1 {
			fmt.Println(line)
		} else {
			fmt.Print(line)
		}

	}

	return nil
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

	switch command {
	case "", "h", "-h", "--h", "/h", "/?", "help", "-help", "--help", "/help":
		printUsage()
		os.Exit(1)
	case "version", "-version", "--version", "/version":
		fmt.Println(version)
		os.Exit(0)
	case "add":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		if value == "" {
			printUsage()
			os.Exit(1)
		}
		err := addValue(file, value, *lastNPtr)
		if err != nil {
			log.Fatal(err)
		}
	case "clear":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		err := clearValues(file)
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
		err := removeValue(file, value)
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		if file == "" {
			printUsage()
			os.Exit(1)
		}
		err := listValues(file, *lastNPtr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		printUsage()
		os.Exit(0)
	}

	return
}
