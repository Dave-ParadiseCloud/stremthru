package endpoint

import (
	"net/http"
	"time"

	"github.com/MunifTanjim/stremthru/internal/config"
	"github.com/MunifTanjim/stremthru/internal/context"
)

type HealthData struct {
	Status string `json:"status"`
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := &HealthData{}
	health.Status = "ok"
	SendResponse(w, 200, health, nil)
}

type HealthDebugDataIP struct {
	Exposed string `json:"exposed"`
	Machine string `json:"machine"`
}

type HealthDebugDataStore struct {
	Default string   `json:"default"`
	Names   []string `json:"names"`
}

type HealthDebugDataUser struct {
	Name  string               `json:"name"`
	Store HealthDebugDataStore `json:"store"`
}

type HealthDebugData struct {
	Time    string               `json:"time"`
	Version string               `json:"version"`
	User    *HealthDebugDataUser `json:"user,omitempty"`
	IP      *HealthDebugDataIP   `json:"ip,omitempty"`
}

func handleHealthDebug(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetRequestContext(r)

	data := &HealthDebugData{
		Time:    time.Now().Format(time.RFC3339),
		Version: config.Version,
	}

	if ctx.IsProxyAuthorized {
		data.User = &HealthDebugDataUser{
			Name: ctx.ProxyAuthUser,
			Store: HealthDebugDataStore{
				Default: config.StoreAuthToken.GetPreferredStore(ctx.ProxyAuthUser),
				Names:   config.StoreAuthToken.ListStores(ctx.ProxyAuthUser),
			},
		}
		exposedIp, _ := config.IP.GetTunnelIP()
		machineIp := config.IP.GetMachineIP()
		data.IP = &HealthDebugDataIP{
			Exposed: exposedIp,
			Machine: machineIp,
		}
	}

	SendResponse(w, 200, data, nil)
}

func AddHealthEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/v0/health", handleHealth)
	mux.HandleFunc("/v0/health/__debug__", Middleware(ProxyAuthContext, StoreContext)(handleHealthDebug))
}
