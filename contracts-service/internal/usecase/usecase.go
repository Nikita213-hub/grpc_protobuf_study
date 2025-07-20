package usecase

import (
	"errors"
	"regexp"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/events"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/repository"
)

type ContractsUsecase interface {
	CreateContract(*models.Contract) error
	GetContract(string) (*models.Contract, error)
	UpdateContract(contractID string, update *models.ContractUpdate) error
}

type contractsUsecase struct {
	contractsRepository repository.ContractsRepository
	publisher           events.ContractEventPublisher
}

func NewContractsUsecase(r repository.ContractsRepository, p events.ContractEventPublisher) ContractsUsecase {
	return &contractsUsecase{
		contractsRepository: r,
		publisher:           p,
	}
}

func (cu *contractsUsecase) CreateContract(contract *models.Contract) error {
	err := cu.contractsRepository.AddContract(contract)
	if err != nil {
		return err
	}
	err = cu.publisher.PublishContractEvent(&events.ContractEvent{
		EventId:   2,
		Type:      events.ContractCreated,
		Contract:  contract,
		OccuredAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (cu *contractsUsecase) GetContract(contractId string) (*models.Contract, error) {
	contract, err := cu.contractsRepository.GetContract(contractId)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (cu *contractsUsecase) UpdateContract(contractID string, update *models.ContractUpdate) error {
	if update.CompanyName != nil && *update.CompanyName == "" {
		return errors.New("company name cannot be empty")
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if update.ContactEmail != nil && !emailRegex.MatchString(*update.ContactEmail) {
		return errors.New("invalid email format")
	}

	if update.MonthlyLimit != nil && *update.MonthlyLimit < 1000 {
		return errors.New("monthly limit must be at least 1000")
	}

	if err := cu.contractsRepository.UpdateContract(contractID, update); err != nil {
		return err
	}

	return nil
}
