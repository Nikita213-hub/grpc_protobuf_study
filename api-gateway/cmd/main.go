package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

//TODO: distribute into different files, Add other endpoints like /contracts, use env for client address

type Handlers struct {
	grpcClient authV2.AuthServiceClient
}

type AuthMiddleware struct {
	authClient authV2.AuthServiceClient
}

func NewAuthMiddleware(authClient authV2.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func NewHandlers(grpcClient authV2.AuthServiceClient) *Handlers {
	return &Handlers{
		grpcClient: grpcClient,
	}
}

func main() {

	authConn, err := grpc.NewClient(
		"localhost:44044",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := authConn.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	authClient := authV2.NewAuthServiceClient(authConn)
	authMiddleware := NewAuthMiddleware(authClient)

	mux := http.NewServeMux()
	handlers := NewHandlers(authClient)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /verify", handlers.VerifyHandler)
	mux.Handle("GET /contracts", authMiddleware.Middleware(ContractsHandler()))
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err)
	}

}

func ContractsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello, contracts"))
	})
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		_, err := am.authClient.GetSession(r.Context(), &authV2.SessionRequest{
			SessionId: sessionID,
		})

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.grpcClient.StartLogin(r.Context(), &authV2.LoginRequest{
		Email: req.Email,
	})

	if err != nil {
		handleGRPCError(w, err)
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

	grpcResp, err := h.grpcClient.VerifyCode(r.Context(), &authV2.VerifyCodeRequest{
		Email: req.Email,
		Code:  req.Code,
	})

	if err != nil {
		handleGRPCError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"session_id": grpcResp.SessionId,
	})
}

func handleGRPCError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unauthenticated:
			http.Error(w, st.Message(), http.StatusUnauthorized)
		case codes.InvalidArgument:
			http.Error(w, st.Message(), http.StatusBadRequest)
		case codes.NotFound:
			http.Error(w, st.Message(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
