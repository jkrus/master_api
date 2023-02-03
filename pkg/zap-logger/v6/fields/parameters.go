package fields

import (
	"go.uber.org/zap"
)

// Parameters ...
//
// It's useful for url.Values and http.Header types.
//
func Parameters(key string, parameters map[string][]string) zap.Field {
	if len(parameters) == 0 {
		return zap.Skip()
	}

	data := make(map[string]interface{})
	for key, values := range parameters {
		var nonemptyValues []string
		for _, value := range values {
			if len(value) != 0 {
				nonemptyValues = append(nonemptyValues, value)
			}
		}

		switch len(nonemptyValues) {
		case 0:
			data[key] = nil
		case 1:
			data[key] = nonemptyValues[0]
		default:
			data[key] = nonemptyValues
		}
	}

	return zap.Reflect(key, data)
}
