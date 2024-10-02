package uuid

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

func GenerateUUIDv7(s string) (string, error) {
	if s == "" {
		guid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return guid.String(), nil
	}

	atTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return "", fmt.Errorf("could not parse time: %w", err)
	}
	guid, err := uuid.NewV7AtTime(atTime)
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}
