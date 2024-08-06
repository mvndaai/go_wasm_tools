package htmltools

import (
	b64 "encoding/base64"
	"html"
)

func Escape(s string) (string, error) {
	return html.EscapeString(s), nil

}

func Unescape(s string) (string, error) {
	return html.UnescapeString(s), nil
}

func B64Encode(s string) (string, error) {
	return b64.StdEncoding.EncodeToString([]byte(s)), nil
}

func B64Decode(s string) (string, error) {
	data, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
