package models

import (
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ContractState int32

const (
	ContractState_UNSPECIFIED ContractState = 0
	ContractState_DRAFT       ContractState = 1
	ContractState_ACTIVE      ContractState = 2
	ContractState_SUSPENDED   ContractState = 3
	ContractState_TERMINATED  ContractState = 4
)

type Contract struct {
	ID             string
	InitiatorID    string
	CompanyName    string
	ContactEmail   string
	MonthlyLimit   float32
	CurrentBalance float32
	State          contractsV1.ContractState
	StartDate      *timestamppb.Timestamp
	EndDate        *timestamppb.Timestamp
	CreatedAt      *timestamppb.Timestamp
	UpdatedAt      *timestamppb.Timestamp
}
