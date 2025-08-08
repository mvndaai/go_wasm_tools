package uuid_test

import (
	"fmt"
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

func TestTimestampFromUUID(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		timezone      string
		format        string
		expected      string
		errorContains string
	}{
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "", expected: "2025-08-08T17:14:26Z"},
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "", expected: "2025-08-08T17:14:26Z"},
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "MST", expected: "2025-08-08T10:14:26-07:00"},
		//{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "Local", expected: "2025-08-08T11:14:26-06:00"}, // Local timezone will vary based on the system running the test
		{name: "uuidv4", in: "af783332-10d5-4314-97ab-3721c842d466", timezone: "", expected: "8083-09-25T06:04:15Z"},
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "UTC", format: "DateTime", expected: "2025-08-08 17:14:26"},
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "UTC", format: "RFC1123", expected: "Fri, 08 Aug 2025 17:14:26 UTC"},
		{name: "uuidv7", in: "01988aad-40c6-7670-8394-811decdfde9c", timezone: "UTC", format: "Kitchen", expected: "5:14PM"},
		{name: "invalid", in: "invalid-uuid", errorContains: "could not parse UUID"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s(tz:%s)", tt.name, tt.timezone), func(t *testing.T) {
			timestamp, err := uuid.TimestampUUIDv7(tt.in + " " + tt.timezone + " " + tt.format)
			if tt.errorContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, timestamp)
			}
		})
	}
}
