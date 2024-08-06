package htmltools_test

import (
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/htmltools"
	"github.com/stretchr/testify/assert"
)

func TestEscapeUnescape(t *testing.T) {
	tests := []struct {
		unescaped string
		escaped   string
	}{
		{unescaped: "<div>", escaped: "&lt;div&gt;"},
		{unescaped: "&'<>\"", escaped: "&amp;&#39;&lt;&gt;&#34;"},
		{unescaped: "\u003C", escaped: "&lt;"},
	}

	for _, tt := range tests {
		t.Run(tt.unescaped, func(t *testing.T) {
			actualEscaped, _ := htmltools.Escape(tt.unescaped)
			assert.Equal(t, tt.escaped, actualEscaped)

			actualUnescaped, _ := htmltools.Unescape(actualEscaped)
			assert.Equal(t, tt.unescaped, actualUnescaped)
		})
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		escaped   string
		unescaped string
	}{
		{escaped: "&lt;div&gt;", unescaped: "<div>"},
		{escaped: "&amp;&#39;&lt;&gt;&#34;", unescaped: "&'<>\""},
		{escaped: "&apos;&quot;", unescaped: "'\""}, // Other versions
		{escaped: "\u003chtml\u003e\r\n\u003c/html\u003e", unescaped: "<html>\r\n</html>"},
	}

	for _, tt := range tests {
		t.Run(tt.unescaped, func(t *testing.T) {
			actualUnescaped, _ := htmltools.Unescape(tt.escaped)
			assert.Equal(t, tt.unescaped, actualUnescaped)
		})
	}
}

func TestB64(t *testing.T) {
	tests := []struct {
		decoded string
		encoded string
	}{
		{decoded: "Hello, World!", encoded: "SGVsbG8sIFdvcmxkIQ=="},
		{decoded: "Hello, 世界!", encoded: "SGVsbG8sIOS4lueVjCE="},
	}

	for _, tt := range tests {
		t.Run(tt.decoded, func(t *testing.T) {
			actualEncoded, _ := htmltools.B64Encode(tt.decoded)
			assert.Equal(t, tt.encoded, actualEncoded)

			actualDecoded, _ := htmltools.B64Decode(actualEncoded)
			assert.Equal(t, tt.decoded, actualDecoded)
		})
	}
}

func TestUrlEncodeDecode(t *testing.T) {
	tests := []struct {
		decoded string
		encoded string
	}{
		{decoded: `{"foo":"bar"}`, encoded: "%7B%22foo%22%3A%22bar%22%7D"},
	}

	for _, tt := range tests {
		t.Run(tt.decoded, func(t *testing.T) {
			actualEncoded, _ := htmltools.URLEncode(tt.decoded)
			assert.Equal(t, tt.encoded, actualEncoded)

			actualDecoded, _ := htmltools.URLDecode(actualEncoded)
			assert.Equal(t, tt.decoded, actualDecoded)
		})
	}
}
