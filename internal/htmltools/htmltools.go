package htmltools

import (
	"html"
)

func Escape(s string) (string, error) {
	return html.EscapeString(s), nil

}

func Unescape(s string) (string, error) {
	return html.UnescapeString(s), nil
}
