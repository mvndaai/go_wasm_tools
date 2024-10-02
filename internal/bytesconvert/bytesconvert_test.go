package bytesconvert_test

import (
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/bytesconvert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name      string
		encoded   string
		decoded   string
		reecnoded string
	}{
		{name: "normal", encoded: "[102 111 111]", decoded: "foo", reecnoded: "[102 111 111]"},
		{name: "forgot []", encoded: "102 111 111", decoded: "foo", reecnoded: "[102 111 111]"},
		{name: "go bytes", encoded: "{0x66 0x6f 0x6f}", decoded: "foo", reecnoded: "[102 111 111]"},
	}

	for _, tt := range tests {
		t.Run(tt.decoded, func(t *testing.T) {
			decoded, err := bytesconvert.ToString(tt.encoded)
			assert.NoError(t, err)
			assert.Equal(t, decoded, tt.decoded)

			reecnoded, err := bytesconvert.FromString(decoded)
			require.NoError(t, err)
			assert.Equal(t, reecnoded, tt.reecnoded)
		})
	}
}
