package grpcserver

import (
	"context"
	"fmt"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/usecase"
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	contractsV1.UnimplementedContractsServiceServer
	usecase usecase.ContractsUsecase
}

func Register(grpcServer *grpc.Server, usecase usecase.ContractsUsecase) {
	contractsV1.RegisterContractsServiceServer(grpcServer, &Server{usecase: usecase})
	reflection.Register(grpcServer)
}

// func (contractsV1.UnimplementedContractsServiceServer) CreateContract(context.Context, *contractsV1.CreateContractRequest) (*contractsV1.ContractOperationResponse, error)
// func (contractsV1.UnimplementedContractsServiceServer) GetContract(context.Context, *contractsV1.GetContractRequest) (*contractsV1.GetContractResponse, error)
// func (contractsV1.UnimplementedContractsServiceServer) UpdateContract(context.Context, *contractsV1.UpdateContractRequest) (*contractsV1.ContractOperationResponse, error)

func (s *Server) CreateContract(ctx context.Context,
	req *contractsV1.CreateContractRequest) (*contractsV1.ContractOperationResponse, error) {
	initiator := req.GetInitiatorId()
	details := req.GetDetails()
	newContract := &models.Contract{
		ID:             uuid.New().String(),
		CompanyName:    details.GetCompanyName(),
		ContactEmail:   details.GetContactEmail(),
		MonthlyLimit:   details.GetMonthlyLimit(),
		StartDate:      details.GetStartDate(),
		EndDate:        details.GetEndDate(),
		CurrentBalance: 0.0,
		State:          *contractsV1.ContractState_CONTRACT_STATE_ACTIVE.Enum(),
		CreatedAt:      timestamppb.Now(),
		UpdatedAt:      timestamppb.Now(),
		InitiatorID:    initiator,
	}
	err := s.usecase.CreateContract(newContract)
	if err != nil {
		return nil, err
	}
	fmt.Println(mapContractToProto(newContract))
	return &contractsV1.ContractOperationResponse{
		Status:   contractsV1.OperationStatus_OPERATION_STATUS_CREATED,
		Message:  "Contract created successfully",
		Contract: mapContractToProto(newContract),
	}, nil
}

func (s *Server) GetContract(ctx context.Context,
	req *contractsV1.GetContractRequest) (*contractsV1.GetContractResponse, error) {
	contractId := req.GetContractId()
	contract, err := s.usecase.GetContract(contractId)
	if err != nil {
		return nil, err
	}
	return &contractsV1.GetContractResponse{
		Status:   contractsV1.OperationStatus_OPERATION_STATUS_SUCCESS,
		Contract: mapContractToProto(contract),
	}, nil
}

func (s *Server) UpdateContract(ctx context.Context,
	req *contractsV1.UpdateContractRequest) (*contractsV1.ContractOperationResponse, error) {
	contractID := req.GetContractId()
	updateDetails := req.GetDetails()

	update := &models.ContractUpdate{}

	if companyName := updateDetails.GetCompanyName(); companyName != nil {
		value := companyName.GetValue()
		update.CompanyName = &value
	}

	if contactEmail := updateDetails.GetContactEmail(); contactEmail != nil {
		value := contactEmail.GetValue()
		update.ContactEmail = &value
	}

	if monthlyLimit := updateDetails.GetMonthlyLimit(); monthlyLimit != nil {
		value := monthlyLimit.GetValue()
		update.MonthlyLimit = &value
	}

	if endDate := updateDetails.GetEndDate(); endDate != nil {
		t := endDate.AsTime()
		update.EndDate = &t
	}

	if err := s.usecase.UpdateContract(contractID, update); err != nil {
		status := contractsV1.OperationStatus_OPERATION_STATUS_BAD_REQUEST
		return &contractsV1.ContractOperationResponse{
			Status:  status,
			Message: fmt.Sprintf("Update failed: %v", err),
		}, nil
	}

	updatedContract, err := s.usecase.GetContract(contractID)
	if err != nil {
		return &contractsV1.ContractOperationResponse{
			Status:  contractsV1.OperationStatus_OPERATION_STATUS_BAD_REQUEST,
			Message: fmt.Sprintf("Failed to fetch updated contract: %v", err),
		}, nil
	}

	return &contractsV1.ContractOperationResponse{
		Status:   contractsV1.OperationStatus_OPERATION_STATUS_SUCCESS,
		Message:  "Contract updated successfully",
		Contract: mapContractToProto(updatedContract),
	}, nil
}

func mapContractToProto(contract *models.Contract) *contractsV1.Contract {

	return &contractsV1.Contract{
		Id:             contract.ID,
		CompanyName:    contract.CompanyName,
		ContactEmail:   contract.ContactEmail,
		MonthlyLimit:   contract.MonthlyLimit,
		CurrentBalance: contract.CurrentBalance,
		State:          contract.State,
		StartDate:      contract.StartDate,
		EndDate:        contract.EndDate,
		CreatedAt:      contract.CreatedAt,
		UpdatedAt:      contract.UpdatedAt,
	}
}
