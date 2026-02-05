package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		allowedOrigins := getAllowedOrigins()
		origin := request.Header.Get("Origin")

		if isOriginAllowed(origin, allowedOrigins) {
			writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func getAllowedOrigins() []string {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		return []string{
			"http://localhost:3000",
			"http://localhost:8080",
		}
	}
	return strings.Split(originsEnv, ",")
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
