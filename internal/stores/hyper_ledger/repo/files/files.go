package files

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/pkg/errors"
)

type FileContractI interface {
	Create(fileUuid string, file dto.File) error
	Update(fileUuid string, file dto.File) error
	GetByUuid(fileUuid string) (*dto.File, error)
}

func NewFileContractI(config *config.Config, logger *zap.Logger, client *client.Network) FileContractI {
	return &fileContract{
		chainCode: client.GetContract(config.HFChaincodeName),
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

// Create adds a new key with value to the world state
func (sc *fileContract) Create(fileUuid string, file dto.File) error {
	_, err := sc.chainCode.SubmitTransaction("CreateFile", fileUuid, file.UserUuid, strconv.Itoa(int(file.TypeId)), string(file.CheckSum), strconv.Itoa(int(file.StatusId)))
	if err != nil {
		return err
	}

	return nil
}

// Update changes the value with key in the world state
func (sc *fileContract) Update(fileUuid string, file dto.File) error {
	submitResult, commit, err := sc.chainCode.SubmitAsync("TransferFile", client.WithArguments(fileUuid, "Mark"))
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
func (sc *fileContract) GetByUuid(fileUuid string) (*dto.File, error) {
	evaluateResult, err := sc.chainCode.EvaluateTransaction("ReadFile", fileUuid)
	if err != nil {
		return nil, err
	}
	result := &dto.File{}
	err = json.Unmarshal(evaluateResult, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
