package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *APIConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf("<html>\n\t<body>\n\t\t<h1>Welcome, Chirpy Admin</h1>\n\t\t<p>Chirpy has been visited %d times!</p>\n\t</body>\n</html>", cfg.fileserverHits.Load()) + "\n",
	))
}
