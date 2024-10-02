package pickle_test

import (
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/pickle"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		decoded string
		encoded string
	}{
		{decoded: "42", encoded: "\x80\x02U\x0242."},
		{decoded: `{"a": 1}`, encoded: "\x80\x02U\b{\"a\": 1}."},
	}

	for _, tt := range tests {
		t.Run(tt.decoded, func(t *testing.T) {
			decoded, err := pickle.Decode(tt.encoded)
			assert.NoError(t, err)
			assert.Equal(t, decoded, tt.decoded)

			encoded, err := pickle.Encode(tt.decoded)
			assert.NoError(t, err)
			assert.Equal(t, encoded, tt.encoded)
		})
	}
}
