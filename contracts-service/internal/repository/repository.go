package repository

import "github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"

type ContractsRepository interface {
	AddContract(*models.Contract) error
	GetContract(contractId string) (*models.Contract, error)
	UpdateContract(contractID string, update *models.ContractUpdate) error
}
