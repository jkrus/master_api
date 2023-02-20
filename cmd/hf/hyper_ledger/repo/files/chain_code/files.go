package chain_code

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type FileContractI interface {
	Create(fileUuid string, file dto.FileINHF) error
	Update(fileUuid string, file dto.FileINHF) error
	GetByUuid(fileUuid string) (*dto.FileINHF, error)
}

func NewFileContractI(logger *zap.Logger, client *client.Network) FileContractI {
	chaincodeName := os.Getenv("CHAINCODE_NAME")
	return &fileContract{
		chainCode: client.GetContract(chaincodeName),
		logger:    logger,
	}
}

// fileContract contract for handling writing and reading from the world state
type fileContract struct {
	logger    *zap.Logger
	chainCode *client.Contract
}

type File struct {
	Uuid         string
	RedactorUuid string
	Type         string
	CheckSum     string
	Status       int
}

type SimpleContract struct {
	contractapi.Contract
}

// Create adds a new key with value to the world state
func (sc *fileContract) Create(fileUuid string, file dto.FileINHF) error {
	_, err := sc.chainCode.SubmitTransaction("CreateAsset", fileUuid, file.RedactorUuid, file.Type, file.CheckSum, string(file.Status))
	if err != nil {
		return err
	}

	return nil
}

// Update changes the value with key in the world state
func (sc *fileContract) Update(fileUuid string, file dto.FileINHF) error {
	submitResult, commit, err := sc.chainCode.SubmitAsync("TransferAsset", client.WithArguments(fileUuid, "Mark"))
	if err != nil {
		return err
	}

	fmt.Println(submitResult)

	if commitStatus, err := commit.Status(); err != nil {
		return err
	} else if !commitStatus.Successful {
		return errors.Newf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code))
	}

	return nil
}

// GetByUuid returns the file at uuid in the world state
func (sc *fileContract) GetByUuid(fileUuid string) (*dto.FileINHF, error) {
	evaluateResult, err := sc.chainCode.EvaluateTransaction("ReadAsset", fileUuid)
	if err != nil {
		return nil, err
	}
	result := &dto.FileINHF{}
	err = json.Unmarshal(evaluateResult, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
