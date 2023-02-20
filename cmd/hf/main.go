package main

import (
	"github.com/jkrus/master_api/cmd/hf/hyper_ledger"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/config"
)

func main() {

	settings := config.GetConfig()

	// Hyper Lager Store
	hfClient, err := hyper_ledger.NewClient(settings)
	if err != nil {
		return
	}
	hfRepo := hyper_ledger.NewHFRepo(hfClient)

	file := dto.FileINHF{}
	err = hfRepo.Files.FileStore.Create("asd", file)
	if err != nil {
		return
	}

}
