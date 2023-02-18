package hyper_ledger

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/hyper_ledger/repo/files/files_i"
)

// FileContractRepo - интерфейс работы смарт контрактом
type FileContractRepo struct {
	Files files_i.FileStoreI
}

// NewFileContractRepo - конструктор интерфейса работы со смарт контрактом
func NewFileContractRepo(logger *zap.Logger) *FileContractRepo {
	return &FileContractRepo{
		Files: files_i.NewFileContract(logger),
	}
}

func (fcr *FileContractRepo) NewContext() {
	fcr.Files.FileStore.
}
