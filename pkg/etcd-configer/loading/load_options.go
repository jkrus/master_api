package loading

import (
	"context"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
)

// Option ...
type Option struct {
	Name  string
	Value []byte
}

//go:generate mockery -name=OptionsGetter -inpkg -case=underscore -testonly

// OptionsGetter returns options by the prefix of their names.
// This prefix will be part of these names.
type OptionsGetter interface {
	GetOptions(ctx context.Context, namePrefix string) ([]Option, error)
}

//go:generate mockery -name=Configuration -inpkg -case=underscore -testonly

// Configuration ...
type Configuration interface {
	SetOption(name string, value []byte) error
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly

// ErrorHandler ...
type ErrorHandler interface {
	HandleError(err error)
}

// LoadingConfiguration ...
type LoadingConfiguration struct { // nolint: golint
	prefix       string
	timeout      time.Duration
	errorHandler ErrorHandler
	skipErrors   bool
}

// LoadingConfigurationOption ...
type LoadingConfigurationOption func(config *LoadingConfiguration) // nolint: golint

// WithPrefix ...
//
// Default: empty.
//
func WithPrefix(prefix string) LoadingConfigurationOption {
	return func(config *LoadingConfiguration) {
		config.prefix = prefix
	}
}

// WithTimeout ...
//
// Default: skipped.
//
func WithTimeout(timeout time.Duration) LoadingConfigurationOption {
	return func(config *LoadingConfiguration) {
		config.timeout = timeout
	}
}

// WithErrorHandler ...
//
// Default: skipped.
//
func WithErrorHandler(errorHandler ErrorHandler) LoadingConfigurationOption {
	return func(config *LoadingConfiguration) {
		config.errorHandler = errorHandler
	}
}

// SkipErrors ...
//
// Default: false.
//
func SkipErrors() LoadingConfigurationOption {
	return func(config *LoadingConfiguration) {
		config.skipErrors = true
	}
}

// LoadOptions ...
func LoadOptions(
	ctx context.Context,
	getter OptionsGetter,
	configuration Configuration,
	options ...LoadingConfigurationOption,
) error {
	var config LoadingConfiguration // default values
	for _, option := range options {
		option(&config)
	}

	if config.timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, config.timeout)
		defer cancel()
	}

	configurationOptions, err := getter.GetOptions(ctx, config.prefix)
	if err != nil {
		err = pkgerrors.WithMessagef(err, "unable to load options with the prefix %s", config.prefix)
		if config.errorHandler != nil {
			config.errorHandler.HandleError(err)
		}
		if !config.skipErrors {
			return err
		}
		return nil
	}

	for _, configurationOption := range configurationOptions {
		name := strings.TrimPrefix(configurationOption.Name, config.prefix)
		if err := configuration.SetOption(name, configurationOption.Value); err != nil {
			err = pkgerrors.WithMessagef(err, "unable to set the option %s", name)
			if config.errorHandler != nil {
				config.errorHandler.HandleError(err)
			}
			if !config.skipErrors {
				return err
			}
		}
	}

	return nil
}
