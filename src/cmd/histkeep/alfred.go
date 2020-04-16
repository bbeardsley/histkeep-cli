package main

import (
	"fmt"
	"regexp"
	"strings"
)

type alfred struct {
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
	globalVars   arrayFlags
}

func (a alfred) list(values []string) {
	fmt.Println("{")

	gvars := buildVariables(a.globalVars, "")
	if len(gvars) > 0 {
		writeVariables(gvars)
		fmt.Println(",")
	}

	fmt.Println("\"items\": [")

	itemCount := 0
	validFormat := a.format.MatchString(a.filter)
	if a.filter == "" || !validFormat {
		for _, item := range a.cannedItems {
			if writeCannedItem(item, itemCount > 0, a.filterFunc) {
				itemCount++
			}
		}
	}

	if len(values) == 0 && itemCount == 0 {
		if a.filter != "" && validFormat {
			values = append(values, a.filter)
		} else if !validFormat {
			for _, item := range a.cannedItems {
				if writeCannedItem(item, itemCount > 0, func(title string) bool { return true }) {
					itemCount++
				}
			}
		}
	}

	for _, line := range values {
		if itemCount > 0 {
			fmt.Println(",")
		}

		fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\"", replacePlaceholder(a.itemTitle, "VALUE", line), replacePlaceholder(a.itemArg, "VALUE", line))
		if a.itemSubtitle != "" {
			fmt.Printf(",\"subtitle\": \"%v\"", replacePlaceholder(a.itemSubtitle, "VALUE", line))
		}
		if a.iconFilename != "" {
			fmt.Print(",\"icon\": { \"path\": \"")
			fmt.Print(replacePlaceholder(a.iconFilename, "VALUE", line))
			fmt.Print("\"} ")
		}
		if a.copyText != "" {
			fmt.Print(",\"text\": { \"copy\": \"")
			fmt.Print(replacePlaceholder(a.copyText, "VALUE", line))
			fmt.Print("\"} ")
		}
		ivars := buildVariables(a.itemVars, line)
		if len(ivars) > 0 {
			fmt.Print(",")
			writeVariables(ivars)
		}

		fmt.Print("}")
		itemCount++
	}
	fmt.Println()
	fmt.Println("]")

	fmt.Println("}")
}

func buildVariables(vars arrayFlags, itemValue string) map[string]string {
	m := make(map[string]string)
	for _, avar := range vars {
		for varName, varValue := range mapNameValuePairs(avar) {
			m[varName] = replacePlaceholder(varValue, "VALUE", itemValue)
		}
	}

	return m
}

func writeVariables(vars map[string]string) {
	fmt.Print("\"variables\": {")

	count := 0
	for varName, varValue := range vars {
		if count > 0 {
			fmt.Print(",")
		}
		fmt.Printf("\"%v\": \"%v\"", varName, varValue)
		count++
	}
	fmt.Print("}")
}

func writeCannedItem(item string, isNotFirst bool, filterFunc func(string) bool) bool {
	pairs := mapNameValuePairs(item)
	if len(pairs) == 0 {
		if filterFunc(item) {
			if isNotFirst {
				fmt.Println(",")
			}
			fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\", \"subtitle\": \"Open %v\"}", item, strings.ToLower(item), strings.ToLower(item))
			return true
		}
	} else {
		title, okTitle := pairs["title"]
		if okTitle && filterFunc(title) {
			arg, okArg := pairs["arg"]
			if !okArg {
				arg = strings.ToLower(title)
			}
			subtitle, okSubtitle := pairs["subtitle"]
			if !okSubtitle {
				subtitle = "Open " + arg
			}
			if isNotFirst {
				fmt.Println(",")
			}
			fmt.Printf("{\"title\": \"%v\",\"arg\": \"%v\", \"subtitle\": \"%v\"}", title, arg, subtitle)
			return true
		}
	}
	return false
}

func replacePlaceholder(input string, placeholder string, replacementValue string) string {
	return strings.Replace(input, "{{"+placeholder+"}}", replacementValue, -1)
}

func mapNameValuePairs(input string) map[string]string {
	m := make(map[string]string)

	results := strings.Split(input, "||")
	if results != nil {
		for i := 0; i < len(results); i++ {
			item := results[i]
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				m[parts[0]] = parts[1]
			}
		}
	}

	return m
}
