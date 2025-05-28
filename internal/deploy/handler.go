package deploy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func HandleDeploy(w http.ResponseWriter, r *http.Request) {
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

	composeDir := os.Getenv("COMPOSE_PROJECT_DIR")
	if composeDir == "" {
		http.Error(w, "COMPOSE_PROJECT_DIR not set", http.StatusInternalServerError)
		log.Println("❌ COMPOSE_PROJECT_DIR not set in env")
		return
	}

	// Run `docker compose pull` and `docker compose up -d --no-deps --build` in that directory
	pullCmd := exec.Command("docker", "compose", "pull")
	pullCmd.Dir = composeDir

	upCmd := exec.Command("docker", "compose", "up", "-d", "--build")
	upCmd.Dir = composeDir

	if output, err := pullCmd.CombinedOutput(); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error pulling images:\n%s", output), http.StatusInternalServerError)
		log.Printf("❌ docker compose pull failed: %s", output)
		return
	}

	if output, err := upCmd.CombinedOutput(); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error restarting services:\n%s", output), http.StatusInternalServerError)
		log.Printf("❌ docker compose up failed: %s", output)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Docker Compose deployment updated successfully"))
	log.Println("✅ Docker Compose deployment updated successfully")
}
