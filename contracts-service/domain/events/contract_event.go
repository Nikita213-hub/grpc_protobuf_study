package events

import (
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"
)

type ContractEventType string

const (
	ContractCreated ContractEventType = "contract_created"
	ContractUpdated ContractEventType = "contract_updated"
)

type ContractEvent struct {
	EventId   int               `json:"event_id"`
	Type      ContractEventType `json:"event_type"`
	Contract  *models.Contract  `json:"contract"`
	OccuredAt time.Time         `json:"occured_at"`
}

type ContractEventPublisher interface {
	PublishContractEvent(event *ContractEvent) error
}
