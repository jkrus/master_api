package hyper_ledger

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/jkrus/blog_copy/hf/hyper_ledger/repo/files/files_i"
)

// HFRepo - интерфейс работы смарт контрактами
type HFRepo struct {
	Files files_i.FileStoreI
}

// NewHFRepo - конструктор интерфейса работы со смарт контрактами
func NewHFRepo(client *client.Network) *HFRepo {
	return &HFRepo{
		Files: files_i.NewFileContract(client),
	}
}
