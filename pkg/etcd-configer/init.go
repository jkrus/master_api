package etcdconfiger

import (
	"context"
	"reflect"
	"time"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/etcd-configer/configuration"
	loading2 "github.com/jkrus/master_api/pkg/etcd-configer/loading"
	"github.com/jkrus/master_api/pkg/etcd-configer/storage"
)

//go:generate mockery -name=Storage -inpkg -case=underscore -testonly

// Storage ...
type Storage interface {
	loading2.OptionsGetter
	loading2.OptionsWatcher
}

// StorageFactory ...
type StorageFactory func() (Storage, error)

// Configuration ...
type Configuration struct {
	environmentPrefix string
	loadingPrefix     string
	dialTimeout       time.Duration
	loadingTimeout    time.Duration
	updateHandler     configuration.UpdateHandler
	logger            *zap.Logger
	retryDelay        time.Duration
	disableStorage    bool
}

// ConfigurationOption ...
type ConfigurationOption func(config *Configuration)

// WithEnvironmentPrefix ...
//
// Default: empty.
//
// Attention! It affects option names only within environment variables
// (but not within the storage).
//
func WithEnvironmentPrefix(prefix string) ConfigurationOption {
	return func(config *Configuration) {
		config.environmentPrefix = prefix
	}
}

// WithLoadingPrefix ...
//
// Default: empty.
//
// Attention! It affects option names only within the storage
// (but not within environment variables).
//
func WithLoadingPrefix(prefix string) ConfigurationOption {
	return func(config *Configuration) {
		config.loadingPrefix = prefix
	}
}

// WithDialTimeout ...
//
// Default: skipped.
//
func WithDialTimeout(timeout time.Duration) ConfigurationOption {
	return func(config *Configuration) {
		config.dialTimeout = timeout
	}
}

// WithLoadingTimeout ...
//
// Default: skipped.
//
func WithLoadingTimeout(timeout time.Duration) ConfigurationOption {
	return func(config *Configuration) {
		config.loadingTimeout = timeout
	}
}

// WithUpdateHandler ...
//
// Default: skipped.
//
func WithUpdateHandler(updateHandler configuration.UpdateHandler) ConfigurationOption {
	return func(config *Configuration) {
		config.updateHandler = updateHandler
	}
}

// WithLogger ...
//
// Default: skipped.
//
func WithLogger(logger *zap.Logger) ConfigurationOption {
	return func(config *Configuration) {
		config.logger = logger
	}
}

// WithRetryDelay ...
//
// Default: 5 s.
//
func WithRetryDelay(retryDelay time.Duration) ConfigurationOption {
	return func(config *Configuration) {
		config.retryDelay = retryDelay
	}
}

// WithoutStorage ...
//
// Default: false.
//
func WithoutStorage(disableStorage bool) ConfigurationOption {
	return func(config *Configuration) {
		config.disableStorage = disableStorage
	}
}

func newConfiguration(options []ConfigurationOption) Configuration {
	// default values
	config := Configuration{
		environmentPrefix: "",                                      // empty
		loadingPrefix:     "",                                      // empty
		dialTimeout:       0,                                       // skipped
		loadingTimeout:    0,                                       // skipped
		updateHandler:     func(name string, value interface{}) {}, // skipped
		logger:            zap.NewNop(),                            // skipped
		retryDelay:        5 * time.Second,
		disableStorage:    false,
	}
	for _, option := range options {
		option(&config)
	}

	return config
}

// Init ...
func Init(
	ctx context.Context,
	dataType reflect.Type,
	endpoint string,
	options ...ConfigurationOption,
) (configuration.StructuralConfiguration, error) {
	config := newConfiguration(options)

	storageFactory := func() (Storage, error) {
		return storage.NewOptionStorage(endpoint, config.dialTimeout)
	}
	return InitFromEnvironmentAndStorage(ctx, dataType, storageFactory, options...)
}

// InitFromEnvironmentAndStorage ...
func InitFromEnvironmentAndStorage(
	ctx context.Context,
	dataType reflect.Type,
	storageFactory StorageFactory,
	options ...ConfigurationOption,
) (configuration.StructuralConfiguration, error) {
	config := newConfiguration(options)

	dataReflection := reflect.New(dataType)
	if err := envconfig.Process(config.environmentPrefix, dataReflection.Interface()); err != nil {
		err = pkgerrors.WithMessage(err, "unable to load options from environment variables")
		return configuration.StructuralConfiguration{}, err
	}

	updateHandler := func(name string, value interface{}) {
		config.logger.With(zap.String("name", name)).Debug("option has been updated")
		config.updateHandler(name, value)
	}

	validator := validator.New()
	dataConfiguration := configuration.NewStructuralConfiguration(
		dataReflection.Interface(),
		// it's defined by the github.com/kelseyhightower/envconfig package
		configuration.WithOptionNameTag("envconfig"),
		// it's defined by the github.com/go-playground/validator package
		configuration.WithValidatorTag("validate"),
		configuration.WithValidatorHandler(validator.Var),
		configuration.WithUpdateHandler(updateHandler),
	)

	if !config.disableStorage {
		go InitFromStorageWithRetrying(ctx, dataConfiguration, storageFactory, options...)
	}

	return dataConfiguration, nil
}

// InitFromStorageWithRetrying ...
func InitFromStorageWithRetrying(
	ctx context.Context,
	dataConfiguration loading2.Configuration,
	storageFactory StorageFactory,
	options ...ConfigurationOption,
) {
	config := newConfiguration(options)

	for {
		InitFromStorageOnce(ctx, dataConfiguration, storageFactory, options...)

		timer := time.NewTimer(config.retryDelay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

// InitFromStorageOnce ...
func InitFromStorageOnce(
	ctx context.Context,
	dataConfiguration loading2.Configuration,
	storageFactory StorageFactory,
	options ...ConfigurationOption,
) {
	config := newConfiguration(options)

	config.logger.Debug("creating the storage")
	storage, err := storageFactory()
	if err != nil {
		config.logger.With(zap.Error(err)).Error("unable to create the storage")
		return
	}

	config.logger.Debug("loading options")
	if err := loading2.LoadOptions(
		ctx,
		storage,
		dataConfiguration,
		loading2.WithPrefix(config.loadingPrefix),
		loading2.WithTimeout(config.loadingTimeout),
	); err != nil {
		config.logger.With(zap.Error(err)).Error("unable to load options")
		return
	}

	config.logger.Debug("watching options")
	watchingLogger := loading2.NewErrorLogger(config.logger, "error on options watching")
	// it will be blocked until the channel got from the storage is closed
	loading2.WatchOptions(
		ctx,
		storage,
		dataConfiguration,
		loading2.WithPrefix(config.loadingPrefix),
		loading2.WithErrorHandler(watchingLogger),
	)
}
