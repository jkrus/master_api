package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// FileSmartContract provides functions for managing an File
type FileSmartContract struct {
	contractapi.Contract
}

// File describes basic details of what makes up a simple file
type File struct {
	Uuid         string
	RedactorUuid string
	Type         string
	CheckSum     string
	Status       int
	History      []History
}

type History struct {
	RedactorUuid string
	Status       int
	UpdatedAt    string
}

// InitLedger adds a base set of files to the ledger
func (s *FileSmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	files := []File{
		{Uuid: "file1", Type: "type1", CheckSum: "CheckSum1", RedactorUuid: "RedactorUuid1", Status: 1, History: []History{{
			RedactorUuid: "RedactorUuid1", Status: 1,
		}}},
		{Uuid: "file2", Type: "type2", CheckSum: "CheckSum2", RedactorUuid: "RedactorUuid2", Status: 2, History: []History{{
			RedactorUuid: "RedactorUuid2", Status: 2,
		}}},
	}

	for _, file := range files {
		fileJSON, err := json.Marshal(file)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(file.Uuid, fileJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateFile issues a new file to the world state with given details.
func (s *FileSmartContract) CreateFile(ctx contractapi.TransactionContextInterface, fileUuid, fileType, checkSum, redactorUuid, updatedAt string, status int) error {
	exists, err := s.FileExists(ctx, fileUuid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the file %s already exists", fileUuid)
	}

	file := File{
		Uuid:         fileUuid,
		Type:         fileType,
		CheckSum:     checkSum,
		RedactorUuid: redactorUuid,
		Status:       status,
		History: []History{{
			RedactorUuid: redactorUuid,
			Status:       status,
			UpdatedAt:    updatedAt,
		}},
	}
	fileJSON, err := json.Marshal(file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fileUuid, fileJSON)
}

// ReadFile returns the file stored in the world state with given id.
func (s *FileSmartContract) ReadFile(ctx contractapi.TransactionContextInterface, fileUuid string) (*File, error) {
	fileJSON, err := ctx.GetStub().GetState(fileUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if fileJSON == nil {
		return nil, fmt.Errorf("the file %s does not exist", fileUuid)
	}

	var file File
	err = json.Unmarshal(fileJSON, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// UpdateFile updates an existing file in the world state with provided parameters.
func (s *FileSmartContract) UpdateFile(ctx contractapi.TransactionContextInterface, fileUuid, redactorUuid, updatedAt string, status int) error {
	exists, err := s.FileExists(ctx, fileUuid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the file %s does not exist", fileUuid)
	}

	file, err := s.ReadFile(ctx, fileUuid)
	if err != nil {
		return fmt.Errorf("the file %s does not read", fileUuid)
	}

	// overwriting original file with new file
	file.RedactorUuid = redactorUuid
	file.Status = status
	file.History = append(file.History, History{
		RedactorUuid: redactorUuid,
		Status:       status,
		UpdatedAt:    updatedAt,
	})

	fileJSON, err := json.Marshal(file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fileUuid, fileJSON)
}

// DeleteFile deletes an given file from the world state.
func (s *FileSmartContract) DeleteFile(ctx contractapi.TransactionContextInterface, fileUuid string) error {
	exists, err := s.FileExists(ctx, fileUuid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the file %s does not exist", fileUuid)
	}

	return ctx.GetStub().DelState(fileUuid)
}

// FileExists returns true when file with given Uuid exists in world state
func (s *FileSmartContract) FileExists(ctx contractapi.TransactionContextInterface, fileUuid string) (bool, error) {
	fileJSON, err := ctx.GetStub().GetState(fileUuid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return fileJSON != nil, nil
}

// GetAllFiles returns all files found in world state
func (s *FileSmartContract) GetAllFiles(ctx contractapi.TransactionContextInterface) ([]*File, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all files in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var files []*File
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var file File
		err = json.Unmarshal(queryResponse.Value, &file)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func main() {
	fileChaincode, err := contractapi.NewChaincode(&FileSmartContract{})
	if err != nil {
		log.Panicf("Error creating file-transfer-basic chaincode: %v", err)
	}

	if err := fileChaincode.Start(); err != nil {
		log.Panicf("Error starting file-transfer-basic chaincode: %v", err)
	}
}
