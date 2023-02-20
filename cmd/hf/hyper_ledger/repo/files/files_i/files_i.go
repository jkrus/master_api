package files_i

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/jkrus/blog_copy/hf/hyper_ledger/repo/files/chain_code"
	"go.uber.org/zap"
)

type FileStoreI struct {
	FileStore chain_code.FileContractI
}

func NewFileContract(logger *zap.Logger, client *client.Network) FileStoreI {
	return FileStoreI{
		FileStore: chain_code.NewFileContractI(logger, client),
	}
}
