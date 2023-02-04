package loading

import (
	"context"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// OptionUpdate ...
type OptionUpdate struct {
	Option

	Error error
}

//go:generate mockery -name=OptionsWatcher -inpkg -case=underscore -testonly

// OptionsWatcher watches options by the prefix of their names.
// This prefix will be part of these names.
type OptionsWatcher interface {
	WatchOptions(ctx context.Context, namePrefix string) chan OptionUpdate
}

// WatchOptions ...
//
// It will be blocked until the channel got from the watcher is closed.
//
func WatchOptions(
	ctx context.Context,
	watcher OptionsWatcher,
	configuration Configuration,
	options ...LoadingConfigurationOption,
) {
	var config LoadingConfiguration // default values
	for _, option := range options {
		option(&config)
	}

	optionChannel := watcher.WatchOptions(ctx, config.prefix)
	for optionUpdate := range optionChannel {
		if err := optionUpdate.Error; err != nil {
			if config.errorHandler != nil {
				config.errorHandler.HandleError(pkgerrors.WithMessagef(
					err,
					"error with watching for options with the prefix %s",
					config.prefix,
				))
			}

			continue
		}

		name := strings.TrimPrefix(optionUpdate.Name, config.prefix)
		if err := configuration.SetOption(name, optionUpdate.Value); err != nil {
			if config.errorHandler != nil {
				err = pkgerrors.WithMessagef(err, "unable to set the option %s", name)
				config.errorHandler.HandleError(err)
			}
		}
	}
}
