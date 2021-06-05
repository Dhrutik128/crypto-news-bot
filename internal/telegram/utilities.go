package telegram

import (
	"fmt"
	"strings"
)

func markdownEscape(s string) string {
	for _, esc := range markdownEscapes {
		if strings.Contains(s, esc) {
			s = strings.Replace(s, esc, fmt.Sprintf("\\%s", esc), -1)
		}
	}
	return s
}
