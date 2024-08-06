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
