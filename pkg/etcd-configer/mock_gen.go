package etcdconfiger

import (
	"github.com/jkrus/master_api/pkg/etcd-configer/loading"
)

//go:generate mockery -name=LoadingConfiguration -inpkg -case=underscore -testonly

// LoadingConfiguration is used only for generating a mock
// for the loading.Configuration type by the mockery tool.
type LoadingConfiguration interface {
	loading.Configuration
}

//go:generate mockery -name=Updater -inpkg -case=underscore -testonly

// Updater is used only for generating a mock
// for the configuration.UpdateHandler type by the mockery tool.
type Updater interface {
	Update(name string, value interface{})
}

//go:generate mockery -name=StorageCreator -inpkg -case=underscore -testonly

// StorageCreator is used only for generating a mock
// for the StorageFactory type by the mockery tool.
type StorageCreator interface {
	CreateStorage() (Storage, error)
}
