package configuration

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"
)

// ValidatorHandler ...
type ValidatorHandler func(field interface{}, tag string) error

// UpdateHandler ...
type UpdateHandler func(name string, value interface{})

// StructuralConfiguration ...
type StructuralConfiguration struct {
	dataLocker     *sync.RWMutex
	dataReflection reflect.Value

	optionNameTag string
	fields        map[string]reflect.StructField

	validatorTag     string
	validatorHandler ValidatorHandler

	updateHandler UpdateHandler
}

// StructuralConfigurationOption ...
type StructuralConfigurationOption func(configuration *StructuralConfiguration)

// WithOptionNameTag ...
//
// Default: "optionName".
//
func WithOptionNameTag(optionNameTag string) StructuralConfigurationOption {
	return func(configuration *StructuralConfiguration) {
		configuration.optionNameTag = optionNameTag
	}
}

// WithValidatorTag ...
//
// Default: "validator".
//
func WithValidatorTag(validatorTag string) StructuralConfigurationOption {
	return func(configuration *StructuralConfiguration) {
		configuration.validatorTag = validatorTag
	}
}

// WithValidatorHandler ...
//
// Default: skipped.
//
func WithValidatorHandler(validatorHandler ValidatorHandler) StructuralConfigurationOption {
	return func(configuration *StructuralConfiguration) {
		configuration.validatorHandler = validatorHandler
	}
}

// WithUpdateHandler ...
//
// Default: skipped.
//
func WithUpdateHandler(updateHandler UpdateHandler) StructuralConfigurationOption {
	return func(configuration *StructuralConfiguration) {
		configuration.updateHandler = updateHandler
	}
}

// NewStructuralConfiguration ...
func NewStructuralConfiguration(
	data interface{},
	options ...StructuralConfigurationOption,
) StructuralConfiguration {
	dataReflection := reflect.ValueOf(data)
	isNotNilPointer := dataReflection.Kind() == reflect.Ptr && !dataReflection.IsNil()
	if !isNotNilPointer || dataReflection.Elem().Kind() != reflect.Struct {
		// correctness depends on a source code at a place of use,
		// and not on an input data or a state of an environment,
		// so it makes no sense to return an error
		const message = "configuration.NewStructuralConfiguration: " +
			"the data should be a not nil pointer to a structure"
		panic(message)
	}

	configuration := StructuralConfiguration{
		dataLocker:     new(sync.RWMutex),
		dataReflection: dataReflection,

		// default values
		optionNameTag:    "optionName",
		validatorTag:     "validator",
		validatorHandler: func(field interface{}, tag string) error { return nil },
		updateHandler:    func(name string, value interface{}) {},
	}
	for _, option := range options {
		option(&configuration)
	}

	configuration.fields = fields(dataReflection.Elem().Type(), configuration.optionNameTag)

	return configuration
}

// SetOption ...
func (configuration StructuralConfiguration) SetOption(name string, value []byte) error {
	field, ok := configuration.fields[name]
	if !ok {
		// ignore not existing options for convenience
		return nil
	}

	// parse the option value
	parsedValue, err := parseValue(value, field.Type)
	if err != nil {
		return pkgerrors.WithMessage(err, "unable to parse the option value")
	}

	// validate the option value
	if tag, ok := field.Tag.Lookup(configuration.validatorTag); ok {
		if err := configuration.validatorHandler(parsedValue, tag); err != nil {
			return pkgerrors.WithMessage(err, "the option has an incorrect value")
		}
	}

	// update the option value
	func() {
		// wrap setting into a closure call for using the defer statement
		// (otherwise the mutex can be lock forever on a panic)

		configuration.dataLocker.Lock()
		defer configuration.dataLocker.Unlock()

		configuration.dataReflection.Elem().FieldByName(field.Name).Set(reflect.ValueOf(parsedValue))
	}()
	configuration.updateHandler(name, parsedValue)

	return nil
}

// CopyData ...
func (configuration StructuralConfiguration) CopyData() interface{} {
	configuration.dataLocker.RLock()
	defer configuration.dataLocker.RUnlock()

	return configuration.dataReflection.Elem().Interface()
}

func optionName(field reflect.StructField, tag string) string {
	name := field.Name
	if tagValue, ok := field.Tag.Lookup(tag); ok {
		name = tagValue
	}

	return name
}

func fields(dataType reflect.Type, tag string) map[string]reflect.StructField {
	count := dataType.NumField()
	if count == 0 {
		return nil
	}

	fields := make(map[string]reflect.StructField, count)
	for index := 0; index < count; index++ {
		field := dataType.Field(index)
		name := optionName(field, tag)
		fields[name] = field
	}

	return fields
}

func parseValue(value []byte, dataType reflect.Type) (interface{}, error) {
	switch dataType {
	case reflect.TypeOf(time.Duration(0)):
		duration, err := time.ParseDuration(string(value))
		if err != nil {
			return nil, pkgerrors.WithMessage(err, "unable to parse the duration")
		}

		return duration, nil
	default:
		valueReflection := reflect.New(dataType)
		if _, err := fmt.Sscanln(string(value), valueReflection.Interface()); err != nil {
			return nil, pkgerrors.WithMessage(err, "unable to scan the value to the target type")
		}

		return valueReflection.Elem().Interface(), nil
	}
}
