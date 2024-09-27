package renders

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/senither/dalamud-plugin-listing/metrics"
	"github.com/senither/dalamud-plugin-listing/state"
)

var (
	lastGeneratedJsonAt int64 = -1
	cachedContent       []byte
)

func RenderJson(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.JsonRoute)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(getRepositoryAsJson()))
}

func getRepositoryAsJson() []byte {
	if lastGeneratedJsonAt == state.GetRepositoriesLastUpdatedAt() {
		return cachedContent
	}

	lastGeneratedJsonAt = state.GetRepositoriesLastUpdatedAt()
	content, err := json.Marshal(state.GetRepositories())
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	cachedContent = content

	return cachedContent
}
