package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nikita213-hub/grpc_protobuf_study/api-gateway/internal/handlers"
	obs "github.com/Nikita213-hub/grpc_protobuf_study/shared/observability"
	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	contractsV1 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/contract"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthMiddleware struct {
	authClient authV2.AuthServiceClient
}

func NewAuthMiddleware(authClient authV2.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func main() {
	ctx := context.Background()
	telemetry, _ := obs.New(ctx, obs.Config{
		ServiceName:    "api-gateway",
		ServiceVersion: "v0.1.0",
		Environment:    "local",
		OtlpEndpoint:   "0.0.0.0:4317",
		SampleRatio:    1.0,
	})
	shutdown := telemetry.Shutdown

	defer shutdown(ctx)

	authConn, err := grpc.NewClient(
		"localhost:44044",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		panic(err)
	}

	contractsConn, err := grpc.NewClient(
		"localhost:44045",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
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
	// Wrap HTTP server with otelhttp via shared helper
	wrapped := obs.WrapHTTP(mux, "http-server").(http.Handler)
	err = http.ListenAndServe("localhost:8080", wrapped)
	if err != nil {
		fmt.Println(err)
	}

}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		sessionID := strings.TrimPrefix(authHeader, "Bearer ")
		if sessionID == "" {
			http.Error(w, "Empty session ID", http.StatusUnauthorized)
			return
		}
		session, err := am.authClient.GetSession(r.Context(), &authV2.SessionRequest{
			SessionId: sessionID,
		})
		fmt.Println(session)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "session", session)
		ctx = context.WithValue(ctx, "userEmail", session.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
