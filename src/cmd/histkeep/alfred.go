package main

import (
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
)

type alfred struct {
	itemTitle          string
	itemSubtitle       string
	itemArg            string
	iconFilename       string
	copyText           string
	itemVars           arrayFlags
	cannedItems        arrayFlags
	filter             string
	filterFunc         func(string) bool
	format             *regexp.Regexp
	globalVars         arrayFlags
	replacePlaceholder func(string, string, string) string
	itemModifiers      arrayFlags
}

func (a alfred) list(values []string) {
	wf := aw.New()
	wf.Args()
	wf.Run(func() {
		for name, value := range buildVariables(a.globalVars, "", a.replacePlaceholder) {
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
			item := wf.NewItem(a.replacePlaceholder(a.itemTitle, "VALUE", line)).
				Arg(a.replacePlaceholder(a.itemArg, "VALUE", line)).
				Valid(true)
			if a.itemSubtitle != "" {
				item.Subtitle(a.replacePlaceholder(a.itemSubtitle, "VALUE", line))
			}
			if a.iconFilename != "" {
				item.Icon(&aw.Icon{
					Value: a.iconFilename,
					Type:  aw.IconTypeImage,
				})
			}
			if a.copyText != "" {
				item.Copytext(a.replacePlaceholder(a.copyText, "VALUE", line))
			}
			for name, value := range buildVariables(a.itemVars, line, a.replacePlaceholder) {
				item.Var(name, value)
			}

			for modifier, modifierValues := range buildModifiers(a.itemModifiers) {
				var mod *aw.Modifier
				switch modifier {
				case "alt":
					mod = item.Alt().Valid(true)
				case "cmd":
					mod = item.Cmd().Valid(true)
				case "ctrl":
					mod = item.Ctrl().Valid(true)
				case "fn":
					mod = item.Fn().Valid(true)
				default:
					mod = nil
				}
				if mod != nil {
					for name, value := range modifierValues {
						switch name {
						case "arg":
							mod.Arg(a.replacePlaceholder(value, "VALUE", line))
						case "icon":
							mod.Icon(&aw.Icon{
								Value: value,
								Type:  aw.IconTypeImage,
							})
						case "subtitle":
							mod.Subtitle(a.replacePlaceholder(value, "VALUE", line))
						case "var":
							for varName, varValue := range mapNameValuePairs(value) {
								mod.Var(varName, a.replacePlaceholder(varValue, "VALUE", line))
							}
						}
					}
				}

			}
		}

		wf.SendFeedback()
	})
}

func buildModifiers(modifierFlags arrayFlags) map[string]map[string]string {
	modifiers := make(map[string]map[string]string)
	for _, modifier := range modifierFlags {
		parts := strings.SplitN(modifier, ":", 3)
		if len(parts) == 3 {
			if modifiers[parts[0]] == nil {
				modifiers[parts[0]] = make(map[string]string)
			}
			modifiers[parts[0]][parts[1]] = parts[2]
		}
	}
	return modifiers
}

func buildVariables(vars arrayFlags, itemValue string, replacePlaceholder func(string, string, string) string) map[string]string {
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
