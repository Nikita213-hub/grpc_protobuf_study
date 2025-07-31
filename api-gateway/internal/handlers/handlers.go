package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/api-gateway/helpers"
	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Handlers struct {
	authClient      authV2.AuthServiceClient
	contractsClient contractsV1.ContractsServiceClient
}

func NewHandlers(authClient authV2.AuthServiceClient, contractsClient contractsV1.ContractsServiceClient) *Handlers {
	return &Handlers{
		authClient:      authClient,
		contractsClient: contractsClient,
	}
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.authClient.StartLogin(r.Context(), &authV2.LoginRequest{
		Email: req.Email,
	})

	if err != nil {
		helpers.HandleGRPCError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": grpcResp.Message,
	})
}

func (h *Handlers) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.authClient.VerifyCode(r.Context(), &authV2.VerifyCodeRequest{
		Email: req.Email,
		Code:  req.Code,
	})

	if err != nil {
		helpers.HandleGRPCError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"session_id": grpcResp.SessionId,
	})
}

func (h *Handlers) CreateContractHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*authV2.Session)
	if !ok {
		http.Error(w, "Session data missing", http.StatusInternalServerError)
		return
	}

	var req struct {
		CompanyName  string  `json:"company_name"`
		ContactEmail string  `json:"contact_email"`
		MonthlyLimit float32 `json:"monthly_limit"`
		StartDate    string  `json:"start_date"`
		EndDate      string  `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startTime, err := helpers.ParseTime(req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	endTime, err := helpers.ParseTime(req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		return
	}

	grpcReq := &contractsV1.CreateContractRequest{
		InitiatorId: session.Email,
		Details: &contractsV1.ContractDetails{
			CompanyName:  req.CompanyName,
			ContactEmail: req.ContactEmail,
			MonthlyLimit: req.MonthlyLimit,
			StartDate:    startTime,
			EndDate:      endTime,
		},
	}

	grpcResp, err := h.contractsClient.CreateContract(r.Context(), grpcReq)
	if err != nil {
		helpers.HandleGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  grpcResp.Status.String(),
		"message": grpcResp.Message,
		"contract": map[string]interface{}{
			"id":              grpcResp.Contract.Id,
			"company_name":    grpcResp.Contract.CompanyName,
			"contact_email":   grpcResp.Contract.ContactEmail,
			"monthly_limit":   grpcResp.Contract.MonthlyLimit,
			"current_balance": grpcResp.Contract.CurrentBalance,
			"state":           grpcResp.Contract.State.String(),
			"start_date":      grpcResp.Contract.StartDate.AsTime().Format(time.RFC3339),
			"end_date":        grpcResp.Contract.EndDate.AsTime().Format(time.RFC3339),
		},
	})
}

func (h *Handlers) UpdateContractHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*authV2.Session)
	if !ok {
		http.Error(w, "Session data missing", http.StatusInternalServerError)
		return
	}

	contractID := strings.TrimPrefix(r.URL.Path, "/contracts/")
	if contractID == "" {
		http.Error(w, "Contract ID required", http.StatusBadRequest)
		return
	}

	var req struct {
		CompanyName  *string  `json:"company_name,omitempty"`
		ContactEmail *string  `json:"contact_email,omitempty"`
		MonthlyLimit *float32 `json:"monthly_limit,omitempty"`
		EndDate      *string  `json:"end_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcReq := &contractsV1.UpdateContractRequest{
		UpdaterId:  session.Email,
		ContractId: contractID,
		Details:    &contractsV1.ContractUpdate{},
	}

	if req.CompanyName != nil {
		grpcReq.Details.CompanyName = &wrapperspb.StringValue{Value: *req.CompanyName}
	}
	if req.ContactEmail != nil {
		grpcReq.Details.ContactEmail = &wrapperspb.StringValue{Value: *req.ContactEmail}
	}
	if req.MonthlyLimit != nil {
		grpcReq.Details.MonthlyLimit = &wrapperspb.FloatValue{Value: *req.MonthlyLimit}
	}
	if req.EndDate != nil {
		endTime, err := helpers.ParseTime(*req.EndDate)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
		grpcReq.Details.EndDate = endTime
	}

	grpcResp, err := h.contractsClient.UpdateContract(r.Context(), grpcReq)
	if err != nil {
		helpers.HandleGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  grpcResp.Status.String(),
		"message": grpcResp.Message,
	})
}

func (h *Handlers) GetContractHandler(w http.ResponseWriter, r *http.Request) {

	contractID := strings.TrimPrefix(r.URL.Path, "/contracts/")
	if contractID == "" {
		http.Error(w, "Contract ID required", http.StatusBadRequest)
		return
	}

	grpcReq := &contractsV1.GetContractRequest{
		ContractId: contractID,
	}

	grpcResp, err := h.contractsClient.GetContract(r.Context(), grpcReq)
	if err != nil {
		helpers.HandleGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": grpcResp.Status.String(),
		"contract": map[string]interface{}{
			"id":              grpcResp.Contract.Id,
			"company_name":    grpcResp.Contract.CompanyName,
			"contact_email":   grpcResp.Contract.ContactEmail,
			"monthly_limit":   grpcResp.Contract.MonthlyLimit,
			"current_balance": grpcResp.Contract.CurrentBalance,
			"state":           grpcResp.Contract.State.String(),
			"start_date":      grpcResp.Contract.StartDate.AsTime().Format(time.RFC3339),
			"end_date":        grpcResp.Contract.EndDate.AsTime().Format(time.RFC3339),
			"created_at":      grpcResp.Contract.CreatedAt.AsTime().Format(time.RFC3339),
			"updated_at":      grpcResp.Contract.UpdatedAt.AsTime().Format(time.RFC3339),
		},
	})
}
