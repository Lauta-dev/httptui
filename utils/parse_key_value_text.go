package utils

import (
	"http_client/const/prefix"
	"strings"
)


// ParseKeyValueText parsea configuración HTTP desde texto key:value.
func ParseKeyValueText(keyValueText string) map[string]string {
	table := make(map[string]string)

	if keyValueText == "" {
		return table
	}

	parts := strings.Split(keyValueText, "\n")

	for _, part := range parts {
		parts := strings.SplitN(strings.TrimSpace(part), ":", 2)

		if len(parts) == 2 {

			if strings.HasPrefix(parts[0], prefix.CommentPrefix) {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			table[key] = value
		}

	}

	return table
}
