package files_i

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/internal/stores/hyper_ledger/repo/files"
)

type FileStoreI struct {
	FileStore files.FileContractI
}

func NewFileContract(config *config.Config, logger *zap.Logger, client *client.Network) FileStoreI {
	return FileStoreI{
		FileStore: files.NewFileContractI(config, logger, client),
	}
}
