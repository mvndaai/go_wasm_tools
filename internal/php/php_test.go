package php_test

import (
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/php"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		decoded string
		encoded string
	}{
		{decoded: "42", encoded: "s:2:\"42\";"},
		{decoded: `{"a": 1}`, encoded: "s:8:\"{\"a\": 1}\";"},
	}

	for _, tt := range tests {
		t.Run(tt.decoded, func(t *testing.T) {
			encoded, err := php.Encode(tt.decoded)
			require.NoError(t, err)
			assert.Equal(t, encoded, tt.encoded)

			decoded, err := php.Decode(tt.encoded)
			assert.NoError(t, err)
			assert.Equal(t, decoded, tt.decoded)
		})
	}
}
