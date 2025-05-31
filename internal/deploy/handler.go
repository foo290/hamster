package deploy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"encoding/json"
)

func RunTransactionDataMigration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	expectedToken := os.Getenv("MIGRATE_TOKEN")
	auth := strings.TrimSpace(r.Header.Get("Authorization"))

	if expectedToken == "" {
		http.Error(w, "MIGRATE_TOKEN not set", http.StatusInternalServerError)
		log.Println("❌ MIGRATE_TOKEN missing in env")
		return
	}

	if auth != "Bearer "+expectedToken {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		log.Printf("❌ Unauthorized access attempt: got token %q\n", auth)
		return
	}
	financeDir := os.Getenv("FINANCE_DIR")
	if financeDir == "" {
		http.Error(w, "FINANCE_DIR not set", http.StatusInternalServerError)
		log.Println("❌ FINANCE_DIR not set in env")
		return
	}
	migrateCmd := exec.Command(financeDir+"/.venv/bin/python", financeDir+"/main.py")
	migrateCmd.Dir = financeDir

	if output, err := migrateCmd.CombinedOutput(); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error running migration:\n%s", output), http.StatusInternalServerError)
		log.Printf("❌ migration failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Transaction data migration completed successfully"))
	log.Println("✅ Transaction data migration completed successfully")
}

func HandleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedToken := os.Getenv("DEPLOY_TOKEN")
	auth := strings.TrimSpace(r.Header.Get("Authorization"))

	if expectedToken == "" {
		http.Error(w, "DEPLOY_TOKEN not set", http.StatusInternalServerError)
		log.Println("❌ DEPLOY_TOKEN missing in env")
		return
	}

	if auth != "Bearer "+expectedToken {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		log.Printf("❌ Unauthorized access attempt: got token %q\n", auth)
		return
	}

	// Parse JSON request body
	var deployRequest struct {
		ComposeDir string `json:"compose_dir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deployRequest); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Printf("❌ Failed to parse request body: %s", err)
		return
	}

	if deployRequest.ComposeDir == "" {
		http.Error(w, "compose_dir is required in request", http.StatusBadRequest)
		log.Println("❌ compose_dir missing in request")
		return
	}

	// Run `docker compose pull` and `docker compose up -d --no-deps --build` in that directory
	composeDir := os.Getenv("COMPOSE_PROJECT_DIR") + "/" + deployRequest.ComposeDir
	pullCmd := exec.Command("docker", "compose", "pull")
	pullCmd.Dir = composeDir

	upCmd := exec.Command("docker", "compose", "up", "-d", "--build")
	upCmd.Dir = composeDir

	if output, err := pullCmd.CombinedOutput(); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error pulling images:\n%s", output), http.StatusInternalServerError)
		log.Printf("❌ docker compose pull failed: error: %s, output: %s", err, output)
		return
	}

	if output, err := upCmd.CombinedOutput(); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error restarting services:\n%s", output), http.StatusInternalServerError)
		log.Printf("❌ docker compose up failed: error: %s, output: %s", err, output)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Docker Compose deployment updated successfully"))
	log.Println("✅ Docker Compose deployment updated successfully")
}
