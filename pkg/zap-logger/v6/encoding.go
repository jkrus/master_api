package zaplogger

import (
	"errors"
	"fmt"
	"strings"
)

// ParseEncoding ...
func ParseEncoding(encoding string) (string, error) {
	switch loweredEncoding := strings.ToLower(encoding); loweredEncoding {
	case "console", "json":
		return loweredEncoding, nil
	case "":
		return "", errors.New("empty encoding")
	default:
		return "", fmt.Errorf("unknown encoding %s", encoding)
	}
}
