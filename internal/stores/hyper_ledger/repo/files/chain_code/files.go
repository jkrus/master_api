package chain_code

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/errors"
)

type FileContractI interface {
	Create(ctx contractapi.TransactionContextInterface, key string, value []byte) error
	Update(ctx contractapi.TransactionContextInterface, key string, value string) error
	Read(ctx contractapi.TransactionContextInterface, key string) (string, error)
}

func NewFileContractI(logger *zap.Logger) FileContractI {
	return &fileContract{
		Contract: &contractapi.Contract{},
		logger:   logger,
	}
}

func ChainCodeStart(ss FileContractI) error {
	fc, ok := interface{}(&fileContract{}).(fileContract)
	if !ok {
		return errors.Newf("Value %v does not implement ContractInterface", ss)
	}

	cc, err := contractapi.NewChaincode(fc)
	if err != nil {
		return errors.Wrap(err, "Can't create chain code")
	}

	fc.chainCode = cc

	if err = fc.chainCode.Start(); err != nil {
		return errors.Wrapf(err, "Can't start chain code")
	}

	return nil
}

// fileContract contract for handling writing and reading from the world state
type fileContract struct {
	*contractapi.Contract
	logger    *zap.Logger
	chainCode *contractapi.ContractChaincode
}

type File struct {
	Uuid     string // Uuid файла
	Data     []byte // Данные файла
	CheckSum string // Контрольная

}

type SimpleContract struct {
	contractapi.Contract
}

// Create adds a new key with value to the world state
func (sc *fileContract) Create(ctx contractapi.TransactionContextInterface, key string, value []byte) error {
	existing, err := ctx.GetStub().GetState(key)

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create world state pair with key %s. Already exists", key)
	}

	err = ctx.GetStub().PutState(key, value)
	st:= ctx.GetStub()
	st.
	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// Update changes the value with key in the world state
func (sc *fileContract) Update(ctx contractapi.TransactionContextInterface, key string, value string) error {
	existing, err := ctx.GetStub().GetState(key)

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	if existing == nil {
		return fmt.Errorf("Cannot update world state pair with key %s. Does not exist", key)
	}

	err = ctx.GetStub().PutState(key, []byte(value))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// Read returns the value at key in the world state
func (sc *fileContract) Read(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	existing, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", errors.New("Unable to interact with world state")
	}

	if existing == nil {
		return "", fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
	}

	return string(existing), nil
}
