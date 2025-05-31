package main

import (
	"log"
	"net/http"
	"os"

	"github.com/foo290/hamster/internal/deploy"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found (skipping)")
	}

	http.HandleFunc("/deploy", deploy.HandleDeploy)
	http.HandleFunc("/migrate-transaction", deploy.RunTransactionDataMigration)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 hamster running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
