package fields

import (
	"crypto/rand"
	"io"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// HashConfig ...
type HashConfig struct {
	randomSource io.Reader
}

// HashOption ...
type HashOption func(config *HashConfig)

// WithRandomSource ...
//
// Default: crypto/rand.Reader.
//
func WithRandomSource(randomSource io.Reader) HashOption {
	return func(config *HashConfig) {
		config.randomSource = randomSource
	}
}

// Hash ...
func Hash(key string, options ...HashOption) zap.Field {
	// default values
	config := HashConfig{
		randomSource: rand.Reader,
	}
	for _, option := range options {
		option(&config)
	}

	hash, err := uuid.NewRandomFromReader(config.randomSource)
	if err != nil {
		err = errors.WithMessage(err, "unable to create an Uuid")
		return MarkedError(key, err)
	}

	return zap.String(key, hash.String())
}
