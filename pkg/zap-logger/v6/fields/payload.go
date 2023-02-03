package fields

import (
	"encoding/json"
	"fmt"
	"go/token"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockery --name=LoggedPayloadData --inpackage --case=underscore --testonly

// LoggedPayloadData ...
type LoggedPayloadData interface {
	LoggedPayloadData() (interface{}, error)
}

// PayloadConfig ...
type PayloadConfig struct {
	SoftLimit     int
	HardLimit     int
	QuantityLimit int
	Placeholder   string
	Suffixes      []string
}

// PayloadOption ...
type PayloadOption func(config *PayloadConfig)

// WithSoftLimit ...
//
// Default: 5 x 1024 (5 KiB).
//
func WithSoftLimit(limit int) PayloadOption {
	return func(config *PayloadConfig) {
		config.SoftLimit = limit
	}
}

// WithHardLimit ...
//
// Default: 5 x 1024 x 1024 (5 MiB).
//
func WithHardLimit(limit int) PayloadOption {
	return func(config *PayloadConfig) {
		config.HardLimit = limit
	}
}

// WithQuantityLimit ...
//
// Default: 10 000.
//
func WithQuantityLimit(limit int) PayloadOption {
	return func(config *PayloadConfig) {
		config.QuantityLimit = limit
	}
}

// WithPlaceholder ...
//
// Default: "[DELETED]".
//
func WithPlaceholder(placeholder string) PayloadOption {
	return func(config *PayloadConfig) {
		config.Placeholder = placeholder
	}
}

// WithSuffix ...
//
// Default: without a suffix.
//
func WithSuffix(suffixes ...string) PayloadOption {
	return func(config *PayloadConfig) {
		config.Suffixes = suffixes
	}
}

// WithRequestURLSuffix ...
//
// It uses the path from request URL as a suffix.
//
func WithRequestURLSuffix(request *http.Request) PayloadOption {
	return WithSuffix(MakeRequestURLSuffix(request))
}

// WithConfig ...
func WithConfig(anotherConfig PayloadConfig) PayloadOption {
	return func(config *PayloadConfig) {
		*config = anotherConfig
	}
}

// ...
var (
	DefaultPayloadConfig = PayloadConfig{
		SoftLimit:     5 * 1024,        /* 5 KiB */
		HardLimit:     5 * 1024 * 1024, /* 5 MiB */
		QuantityLimit: 10000,
		Placeholder:   "[DELETED]",
		Suffixes:      nil,
	}
)

// Payload ...
//
// It accepts byte slices, strings (both should be JSON) and structures (will be convert to JSON).
//
func Payload(key string, rawValue interface{}, options ...PayloadOption) (fields []zap.Field) {
	config := DefaultPayloadConfig
	for _, option := range options {
		option(&config)
	}

	defer func() {
		for index := range fields {
			// two skipped stack frames for the current closure and its parent function
			fields[index] = AddSuffix(fields[index], 2, config.Suffixes...)
		}
	}()

	if value, ok := rawValue.(LoggedPayloadData); ok {
		data, err := value.LoggedPayloadData()
		if err != nil {
			err = errors.WithMessage(err, "unable to get the logged payload data")
			return []zap.Field{MarkedError(key, err)}
		}

		return []zap.Field{zap.Any(key, data)}
	}

	// convert the raw value to bytes
	var text []byte
	switch value := rawValue.(type) {
	case []byte:
		text = value
	case string:
		text = []byte(value)
	default:
		var err error
		text, err = json.Marshal(value)
		if err != nil {
			err = errors.WithMessage(err, "unable to marshal the value to JSON")
			return []zap.Field{MarkedError(key, err)}
		}
	}

	// try to convert the payload to an inner representation (map[string]interface{})...
	var value interface{}
	var restrictedPaths []string
	if err := json.Unmarshal(text, &value); err == nil {
		// ...if it was successful, restrict the size of fields
		value, restrictedPaths = restrictValueSize("root", value, config)
		// marshaling of the result of unmarshalling, so error can't occur
		text, _ = json.Marshal(value) // nolint: errcheck, gosec
	}

	// construct a main field
	if len(text) <= config.HardLimit {
		fields = []zap.Field{zap.ByteString(key, text)}
	} else {
		fields = []zap.Field{zap.Reflect(key, nil)}
		restrictedPaths = []string{"root"}
	}

	// construct an auxiliary field
	if len(restrictedPaths) != 0 {
		restrictedPathsFieldSuffix := fmt.Sprintf("restricted%spaths", suffixSeparator)
		restrictedPathsField :=
			AddSuffix(zap.Strings(key, restrictedPaths), 0, restrictedPathsFieldSuffix)
		fields = append(fields, restrictedPathsField)
	}

	return fields
}

func restrictValueSize(rawValuePath string, rawValue interface{}, config PayloadConfig) (
	restrictedValue interface{},
	restrictedValuePaths []string,
) {
	switch value := rawValue.(type) {
	case string:
		if len(value) > config.SoftLimit {
			return config.Placeholder, []string{rawValuePath}
		}

		return value, nil
	case []interface{}:
		if len(value) > config.QuantityLimit {
			value = value[:config.QuantityLimit]
			restrictedValuePaths = []string{rawValuePath}
		}

		var restrictedSlice []interface{}
		for childIndex, childValue := range value {
			restrictedChildValuePath := fmt.Sprintf("%s[%d]", rawValuePath, childIndex)
			restrictedChildValue, restrictedChildValuePaths :=
				restrictValueSize(restrictedChildValuePath, childValue, config)

			restrictedSlice = append(restrictedSlice, restrictedChildValue)
			restrictedValuePaths = append(restrictedValuePaths, restrictedChildValuePaths...)
		}

		return restrictedSlice, restrictedValuePaths
	case map[string]interface{}:
		restrictValue := make(map[string]interface{})
		for childKey, childValue := range value {
			if len(restrictValue) == config.QuantityLimit {
				restrictedValuePaths = append(restrictedValuePaths, rawValuePath)
				break
			}

			var restrictedChildValuePath string
			// it would be more correct to check for a valid JavaScript identifier, but I used
			// what I had on hand so as not to over-complicate the code
			if token.IsIdentifier(childKey) {
				restrictedChildValuePath = fmt.Sprintf("%s.%s", rawValuePath, childKey)
			} else {
				restrictedChildValuePath = fmt.Sprintf("%s[%q]", rawValuePath, childKey)
			}

			restrictedChildValue, restrictedChildValuePaths :=
				restrictValueSize(restrictedChildValuePath, childValue, config)

			restrictValue[childKey] = restrictedChildValue
			restrictedValuePaths = append(restrictedValuePaths, restrictedChildValuePaths...)
		}

		return restrictValue, restrictedValuePaths
	default:
		return value, nil
	}
}
