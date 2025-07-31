package helpers

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(w http.ResponseWriter, err error) {
	fmt.Println(err)
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
