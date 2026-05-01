package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	handler "esimply/internal/delivery/http"
	mysqlrepo "esimply/internal/repository/mysql"
	"esimply/internal/usecase"
	"esimply/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/esimply?parseTime=true&charset=utf8mb4"
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	if err := database.SeedIfEmpty(db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	repo := mysqlrepo.NewCountryRepository(db)
	uc := usecase.NewCountryUsecase(repo)
	h := handler.NewHandler(uc)
	router := handler.NewRouter(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
