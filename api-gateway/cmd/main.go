package main

import (
	"fmt"
	"net/http"

	"github.com/Nikita213-hub/grpc_protobuf_study/api-gateway/internal/handlers"
	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//TODO: distribute into different files, Add other endpoints like /contracts, use env for client address

type AuthMiddleware struct {
	authClient authV2.AuthServiceClient
}

func NewAuthMiddleware(authClient authV2.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
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

	contractsConn, err := grpc.NewClient(
		"localhost:44045",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := authConn.Close(); err != nil {
			fmt.Println(err)
		}
		if err := contractsConn.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	authClient := authV2.NewAuthServiceClient(authConn)
	authMiddleware := NewAuthMiddleware(authClient)

	contractsClient := contractsV1.NewContractsServiceClient(contractsConn)

	mux := http.NewServeMux()
	handlers := handlers.NewHandlers(authClient, contractsClient)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /verify", handlers.VerifyHandler)
	mux.Handle("POST /contracts", authMiddleware.Middleware(http.HandlerFunc(handlers.CreateContractHandler)))
	mux.Handle("PATCH /contracts/{id}", authMiddleware.Middleware(http.HandlerFunc(handlers.UpdateContractHandler)))
	mux.Handle("GET /contracts/{id}", authMiddleware.Middleware(http.HandlerFunc(handlers.GetContractHandler)))
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err)
	}

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
