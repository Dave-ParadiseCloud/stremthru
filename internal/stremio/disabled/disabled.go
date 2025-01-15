package stremio_disabled

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MunifTanjim/stremthru/internal/shared"
)

func handleManifest(w http.ResponseWriter, r *http.Request) {
	if !shared.IsMethod(r, http.MethodGet) {
		shared.ErrorMethodNotAllowed(r).Send(w)
		return
	}

	manifest, err := GetDisabledManifest(r.PathValue("manifestUrl"))
	if err != nil {
		shared.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(manifest); err != nil {
		log.Printf("failed to encode json %v\n", err)
	}
}

func handleConfigure(w http.ResponseWriter, r *http.Request) {
	if !shared.IsMethod(r, http.MethodGet) {
		shared.ErrorMethodNotAllowed(r).Send(w)
		return
	}

	http.Redirect(w, r, "/stremio/sidekick/?addon_operation=manage&try_load_addons=1", http.StatusFound)
}

func AddStremioDisabledEndpoints(mux *http.ServeMux) {
	withCors := shared.Middleware(shared.EnableCORS)

	mux.HandleFunc("/stremio/disabled/{manifestUrl}/manifest.json", withCors(handleManifest))
	mux.HandleFunc("/stremio/disabled/{manifestUrl}/configure", handleConfigure)
}
