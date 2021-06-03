package tramsfunc

import (
	"fmt"
	"net/http"
	"os"
)

// API is the entry point for GCP Functions.
func API(w http.ResponseWriter, r *http.Request) {
	var apiKey string
	if apiKey = os.Getenv("TFGM_API_KEY"); apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "TFGM API Key is not set")
		return
	}
	fmt.Fprintf(w, "Hello, World! API KEY = %s\n", apiKey)
}
