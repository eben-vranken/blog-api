package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eben-vranken/blog-api/internal/database"
	"github.com/eben-vranken/blog-api/internal/handlers"
	"github.com/eben-vranken/blog-api/internal/repository"
	"github.com/joho/godotenv"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.New(databaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	userRepository := repository.CreateUserRepository(db)
	userHandler := handlers.CreateUserHandler(&userRepository)

	http.HandleFunc("GET /health", loggingMiddleware(healthCheck))

	http.HandleFunc("GET /user", loggingMiddleware(userHandler.GetAll))
	http.HandleFunc("POST /user", loggingMiddleware(userHandler.Create))
	http.HandleFunc("GET /user/{id}", loggingMiddleware(userHandler.GetSpecific))
	http.HandleFunc("DELETE /user/{id}", loggingMiddleware(userHandler.Delete))

	log.Print("Listening on port 8080...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Everything up and running!"))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print(req.URL.Path, " Initializing logging middleware")
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, req)
		duration := time.Since(start)
		log.Printf("[%s] %s %s %d", req.Method, req.RequestURI, duration, recorder.statusCode)
	})
}
