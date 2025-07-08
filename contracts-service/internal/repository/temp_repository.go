package repository

import "github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"

type TempContractsRepository struct {
}

func NewTempContractsRepository() ContractsRepository {
	return &TempContractsRepository{}
}

func (tr *TempContractsRepository) AddContract(*models.Contract) error {
	return nil
}

func (tr *TempContractsRepository) GetContract(string) (*models.Contract, error) {
	return nil, nil
}

func (tr *TempContractsRepository) UpdateContract(contractID string, update *models.ContractUpdate) error {
	return nil
}
