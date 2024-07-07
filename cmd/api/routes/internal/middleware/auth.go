// package middleware

// import (
// 	"log"
// 	"net/http"
// )

// func AuthMiddleware(next http.Handler) {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// authHeader := r.Header.Get("Authorization")
// 		// if authHeader == "" {
// 		// 	http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
// 		// 	return
// 		// }

// 		// authTokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		log.Println("Authorization token is valid")
// 		next.ServeHTTP(w, r)
// 	})
// }

package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) {
	return
}
