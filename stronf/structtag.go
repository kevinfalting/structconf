package stronf

import "strings"

func parseStructTag(tag string) map[string]string {
	options := make(map[string]string)
	for _, option := range strings.Split(tag, ",") {
		optionParts := strings.SplitN(option, ":", 2)
		key := optionParts[0]
		var value string
		if len(optionParts) > 1 {
			value = optionParts[1]
		}

		options[key] = value
	}

	return options
}
