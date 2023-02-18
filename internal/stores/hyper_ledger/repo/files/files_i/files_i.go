package files_i

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/hyper_ledger/repo/files/chain_code"
)

type FileStoreI struct {
	FileStore chain_code.FileContractI
}

func NewFileContract(logger *zap.Logger) FileStoreI {
	return FileStoreI{
		FileStore: chain_code.NewFileContractI(logger),
	}
}
