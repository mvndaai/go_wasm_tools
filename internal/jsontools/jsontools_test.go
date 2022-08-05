package jsontools_test

import (
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/jsontools"
	"github.com/stretchr/testify/assert"
)

func TestJSONConversions(t *testing.T) {
	compressed := `{"a":"b"}`
	escaped := `"{\"a\":\"b\"}"`
	pretty := "{\n\t\"a\": \"b\"\n}"

	e, err := jsontools.Escape(compressed)
	assert.Nil(t, err)
	assert.Equal(t, escaped, e)

	u, err := jsontools.Unescape(e)
	assert.Nil(t, err)
	assert.Equal(t, compressed, u)

	p, err := jsontools.Pretty(compressed)
	assert.Nil(t, err)
	assert.Equal(t, pretty, p)

	c, err := jsontools.Compress(p)
	assert.Nil(t, err)
	assert.Equal(t, compressed, c)
}
