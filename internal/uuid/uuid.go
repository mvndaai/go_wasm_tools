package uuid

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

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

func TimestampUUIDv7(uuidAndTimestamp string) (string, error) {
	parts := strings.Split(uuidAndTimestamp, " ")
	uuidStr := parts[0]
	timezone := ""
	if len(parts) > 1 {
		timezone = parts[1]
	}
	format := ""
	if len(parts) > 2 {
		format = strings.Join(parts[2:], " ")
	}
	return timestampFromUUID(uuidStr, timezone, format)
}

func timestampFromUUID(uuidStr, timezone, format string) (string, error) {
	_, err := uuid.FromString(uuidStr)
	if err != nil {
		return "", fmt.Errorf("could not parse UUID: %w", err)
	}
	t, err := extractTimestamp(uuidStr)
	if err != nil {
		return "", fmt.Errorf("could not extract timestamp: %w", err)
	}

	if timezone != "" {
		t = t.UTC()
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("could not load location: %w", err)
	}
	t = t.In(loc)

	return t.Format(timeFormat(format)), nil
}

// https://gist.github.com/edr3x/5d836a5fdd76adbcba0fad0130ed984d
func extractTimestamp(uuid string) (time.Time, error) {
	parts := strings.Split(uuid, "-")
	millisecondsStr := parts[0] + parts[1]
	milliseconds, err := strconv.ParseInt(millisecondsStr, 16, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestampSeconds := milliseconds / 1000
	return time.Unix(timestampSeconds, 0), nil
}

func timeFormat(fmt string) string {
	switch fmt {
	case "Layout":
		return time.Layout
	case "ANSIC":
		return time.ANSIC
	case "UnixDate":
		return time.UnixDate
	case "RubyDate":
		return time.RubyDate
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC850":
		return time.RFC850
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "", "RFC3339":
		return time.RFC3339
	case "RFC3339Nano":
		return time.RFC3339Nano
	case "Kitchen":
		return time.Kitchen
	case "Stamp":
		return time.Stamp
	case "StampMilli":
		return time.StampMilli
	case "StampMicro":
		return time.StampMicro
	case "StampNano":
		return time.StampNano
	case "DateTime":
		return time.DateTime
	case "DateOnly":
		return time.DateOnly
	case "TimeOnly":
		return time.TimeOnly
	default:
		return fmt // Return the format as is if not recognized.
	}
}
