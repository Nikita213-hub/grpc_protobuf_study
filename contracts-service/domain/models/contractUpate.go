package models

import "time"

type ContractUpdate struct {
	CompanyName  *string
	ContactEmail *string
	MonthlyLimit *float32
	EndDate      *time.Time
}
