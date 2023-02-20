package hyper_ledger

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/internal/stores/hyper_ledger/repo/files/files_i"
)

// HFRepo - интерфейс работы смарт контрактами
type HFRepo struct {
	Files files_i.FileStoreI
}

// NewHFRepo - конструктор интерфейса работы со смарт контрактами
func NewHFRepo(config *config.Config, logger *zap.Logger, client *client.Network) *HFRepo {
	return &HFRepo{
		Files: files_i.NewFileContract(config, logger, client),
	}
}
