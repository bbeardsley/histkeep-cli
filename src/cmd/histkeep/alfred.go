package main

import (
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
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
	wf := aw.New()
	wf.Args()
	wf.Run(func() {
		for name, value := range buildVariables(a.globalVars, "") {
			wf.Var(name, value)
		}

		hasCanned := false
		validFormat := a.format.MatchString(a.filter)
		if a.filter == "" || !validFormat {
			for _, item := range a.cannedItems {
				if writeCannedItem(wf, item, a.filterFunc) {
					hasCanned = true
				}
			}
		}

		if len(values) == 0 && !hasCanned {
			if a.filter != "" && validFormat {
				values = append(values, a.filter)
			} else if !validFormat {
				for _, item := range a.cannedItems {
					writeCannedItem(wf, item, func(title string) bool { return true })
				}
			}
		}

		for _, line := range values {
			item := wf.NewItem(replacePlaceholder(a.itemTitle, "VALUE", line)).
				Arg(replacePlaceholder(a.itemArg, "VALUE", line)).
				Valid(true)
			if a.itemSubtitle != "" {
				item.Subtitle(replacePlaceholder(a.itemSubtitle, "VALUE", line))
			}
			if a.iconFilename != "" {
				item.Icon(&aw.Icon{
					Value: a.iconFilename,
					Type:  aw.IconTypeImage,
				})
			}
			if a.copyText != "" {
				item.Copytext(replacePlaceholder(a.copyText, "VALUE", line))
			}
			for name, value := range buildVariables(a.itemVars, line) {
				item.Var(name, value)
			}
		}

		wf.SendFeedback()
	})
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

func writeCannedItem(wf *aw.Workflow, cannedItem string, filterFunc func(string) bool) bool {
	pairs := mapNameValuePairs(cannedItem)
	if len(pairs) == 0 {
		if filterFunc(cannedItem) {
			wf.NewItem(cannedItem).Arg(strings.ToLower(cannedItem)).Subtitle(strings.ToLower(cannedItem)).Valid(true)
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
			wf.NewItem(title).Arg(arg).Subtitle(subtitle).Valid(true)
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
