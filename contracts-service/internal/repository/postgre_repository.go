package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/models"
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postgreRepository struct {
	dbConn *sql.DB
}

func NewPostgreRepository(conn *sql.DB) ContractsRepository {
	return &postgreRepository{dbConn: conn}
}

func (pr *postgreRepository) AddContract(contract *models.Contract) error {
	var status string
	switch contract.State {
	case contractsV1.ContractState_CONTRACT_STATE_ACTIVE:
		status = "active"
	case contractsV1.ContractState_CONTRACT_STATE_DRAFT:
		status = "expired"
	case contractsV1.ContractState_CONTRACT_STATE_SUSPENDED:
		status = "blocked"
	default:
		return fmt.Errorf("invalid contract state: %v", contract.State)
	}

	startDate := contract.StartDate.AsTime()
	endDate := contract.EndDate.AsTime()
	createdAt := contract.CreatedAt.AsTime()
	updatedAt := contract.UpdatedAt.AsTime()

	query := `
	INSERT INTO contracts (
		id, 
		company_name, 
		contact_email, 
		monthly_limit, 
		current_balance, 
		start_date, 
		end_date, 
		status, 
		created_at, 
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := pr.dbConn.Exec(query,
		contract.ID,
		contract.CompanyName,
		contract.ContactEmail,
		contract.MonthlyLimit,
		contract.CurrentBalance,
		startDate,
		endDate,
		status,
		createdAt,
		updatedAt,
	)

	return err
}

func (pr *postgreRepository) GetContract(contractId string) (*models.Contract, error) {
	query := `
        SELECT 
            id, 
            company_name, 
            contact_email, 
            monthly_limit, 
            current_balance, 
            start_date, 
            end_date, 
            status, 
            created_at, 
            updated_at
        FROM contracts 
        WHERE id = $1
    `

	row := pr.dbConn.QueryRow(query, contractId)

	var (
		id             string
		companyName    string
		contactEmail   string
		monthlyLimit   float32
		currentBalance float32
		startDate      time.Time
		endDate        time.Time
		status         string
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := row.Scan(
		&id,
		&companyName,
		&contactEmail,
		&monthlyLimit,
		&currentBalance,
		&startDate,
		&endDate,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found with id: %s", contractId)
		}
		return nil, fmt.Errorf("failed to scan contract: %w", err)
	}

	var state contractsV1.ContractState
	switch status {
	case "active":
		state = *contractsV1.ContractState_CONTRACT_STATE_ACTIVE.Enum()
	case "blocked":
		state = *contractsV1.ContractState_CONTRACT_STATE_SUSPENDED.Enum()
	case "expired":
		state = *contractsV1.ContractState_CONTRACT_STATE_TERMINATED.Enum()
	default:
		state = *contractsV1.ContractState_CONTRACT_STATE_UNSPECIFIED.Enum()
	}

	contract := &models.Contract{
		ID:             id,
		CompanyName:    companyName,
		ContactEmail:   contactEmail,
		MonthlyLimit:   monthlyLimit,
		CurrentBalance: currentBalance,
		State:          state,
		StartDate:      timestamppb.New(startDate),
		EndDate:        timestamppb.New(endDate),
		CreatedAt:      timestamppb.New(createdAt),
		UpdatedAt:      timestamppb.New(updatedAt),
	}

	return contract, nil
}

func (pr *postgreRepository) UpdateContract(contractID string, update *models.ContractUpdate) error {
	setClauses := []string{"updated_at = $1"}
	params := []interface{}{
		time.Now(),
	}
	index := 2

	if update.CompanyName != nil {
		setClauses = append(setClauses, fmt.Sprintf("company_name = $%d", index))
		params = append(params, *update.CompanyName)
		index++
	}

	if update.ContactEmail != nil {
		setClauses = append(setClauses, fmt.Sprintf("contact_email = $%d", index))
		params = append(params, *update.ContactEmail)
		index++
	}

	if update.MonthlyLimit != nil {
		setClauses = append(setClauses, fmt.Sprintf("monthly_limit = $%d", index))
		params = append(params, *update.MonthlyLimit)
		index++
	}

	if update.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", index))
		params = append(params, *update.EndDate)
		index++
	}

	if len(setClauses) == 1 {
		return nil
	}

	params = append(params, contractID)
	whereClause := fmt.Sprintf("id = $%d", index)

	query := fmt.Sprintf("UPDATE contracts SET %s WHERE %s",
		strings.Join(setClauses, ", "),
		whereClause,
	)

	_, err := pr.dbConn.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
