package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

// SetupSwaggerRoutes настраивает маршруты для Swagger документации
func SetupSwaggerRoutes(r *mux.Router) {
	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger-doc.json"),
	))
	
	// Swagger JSON
	r.HandleFunc("/swagger-doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		_, err := w.Write([]byte(SwaggerDoc))
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	})
}
