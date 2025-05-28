package deploy

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func HandleDeploy(w http.ResponseWriter, r *http.Request) {
	expectedToken := os.Getenv("DEPLOY_TOKEN")
	auth := r.Header.Get("Authorization")

	if auth != "Bearer "+expectedToken {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	deployment := os.Getenv("K8S_DEPLOYMENT_NAME")
	if deployment == "" {
		http.Error(w, "K8S_DEPLOYMENT_NAME not set", http.StatusInternalServerError)
		return
	}

	namespace := os.Getenv("K8S_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	cmd := exec.Command("kubectl", "rollout", "restart", fmt.Sprintf("deployment/%s", deployment), "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("❌ Error: %s", output), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("✅ Restarted deployment: %s", output)))
}
