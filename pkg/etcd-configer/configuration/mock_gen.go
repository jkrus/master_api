package configuration

//go:generate mockery -name=Validator -inpkg -case=underscore -testonly

// Validator is used only for generating a mock
// for the ValidatorHandler type by the mockery tool.
type Validator interface {
	Validate(field interface{}, tag string) error
}

//go:generate mockery -name=Updater -inpkg -case=underscore -testonly

// Updater is used only for generating a mock
// for the UpdateHandler type by the mockery tool.
type Updater interface {
	Update(name string, value interface{})
}
