package main

import (
	"fmt"
	"regexp"
	"strings"
)

type alfred struct {
	values       []string
	itemTitle    string
	itemSubtitle string
	itemArg      string
	iconFilename string
	copyText     string
	itemVars     arrayFlags
	cannedItems  arrayFlags
	filter       string
	filterFunc   func(string) bool
	format       *regexp.Regexp
}

func (a alfred) list() {
	fmt.Println("{\"items\": [")

	itemCount := 0
	validFormat := a.format.MatchString(a.filter)
	if a.filter == "" || !validFormat {
		for _, item := range a.cannedItems {
			if a.filterFunc(item) {
				if itemCount > 0 {
					fmt.Println(",")
				}
				fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\", \"subtitle\": \"Open %v\"}", item, strings.ToLower(item), strings.ToLower(item))
				itemCount = itemCount + 1
			}
		}
	}

	if len(a.values) == 0 && itemCount == 0 {
		if a.filter != "" && validFormat {
			a.values = append(a.values, a.filter)
		} else if !validFormat {
			for _, item := range a.cannedItems {
				if itemCount > 0 {
					fmt.Println(",")
				}
				fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\", \"subtitle\": \"Open %v\"}", item, strings.ToLower(item), strings.ToLower(item))
				itemCount = itemCount + 1
			}
		}
	}

	for _, line := range a.values {
		if itemCount > 0 {
			fmt.Println(",")
		}

		fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\"", replacePlaceholder(a.itemTitle, "VALUE", line), replacePlaceholder(a.itemArg, "VALUE", line))
		if a.itemSubtitle != "" {
			fmt.Printf(", \"subtitle\": \"%v\"", replacePlaceholder(a.itemSubtitle, "VALUE", line))
		}
		if a.iconFilename != "" {
			fmt.Print(", \"icon\": { \"path\": \"")
			fmt.Print(replacePlaceholder(a.iconFilename, "VALUE", line))
			fmt.Print("\"} ")
		}
		if a.copyText != "" {
			fmt.Print(", \"text\": { \"copy\": \"")
			fmt.Print(replacePlaceholder(a.copyText, "VALUE", line))
			fmt.Print("\"} ")
		}
		if len(a.itemVars) > 0 {
			fmt.Print(", \"variables\": {")
			for i, avar := range a.itemVars {
				if i != 0 {
					fmt.Print(",")
				}
				parts := strings.Split(avar, "=")
				if len(parts) == 2 {
					fmt.Printf("\"%v\": \"%v\"", replacePlaceholder(parts[0], "VALUE", line), replacePlaceholder(parts[1], "VALUE", line))
				}
			}
			fmt.Println("}")
		}
		fmt.Print("}")
		itemCount = itemCount + 1
	}
	fmt.Println()
	fmt.Println("]}")
}

func replacePlaceholder(input string, placeholder string, replacementValue string) string {
	return strings.Replace(input, "{{"+placeholder+"}}", replacementValue, -1)
}
