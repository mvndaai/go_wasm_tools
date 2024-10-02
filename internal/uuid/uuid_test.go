package uuid_test

import (
	"strings"
	"testing"

	"github.com/mvndaai/go_wasm_tools/internal/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPrefixes(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		prefix string
	}{
		{name: "from time", in: "2021-01-01T00:00:00Z", prefix: "0176bb3e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := uuid.GenerateUUIDv7(tt.in)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(out, tt.prefix), out)
		})
	}

}
