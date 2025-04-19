package renders

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/senither/dalamud-plugin-listing/metrics"
	"github.com/senither/dalamud-plugin-listing/state"
)

var (
	fileHash *string
)

func RenderHtml(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.HtmlRoute)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	content, err := os.ReadFile("./views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	if fileHash == nil {
		hash := createSHA1Hash(string(content))
		fileHash = &hash
	}

	fmt.Fprint(w, renderTemplateStrings(string(content)))
}

func renderTemplateStrings(template string) string {
	template = strings.Replace(template, "@file-hash", *fileHash, 1)
	template = strings.ReplaceAll(template, "@state-url-size", strconv.Itoa(state.GetUrlsSize()))
	template = strings.ReplaceAll(template, "@state-repo-size", strconv.Itoa(state.GetRepositoriesSize()))
	template = strings.ReplaceAll(template, "@state-internal-size", strconv.Itoa(state.GetInternalPluginSize()))
	template = strings.ReplaceAll(template, "@state-senither-size", strconv.Itoa(state.GetSenitherPluginSize()))

	return template
}

func createSHA1Hash(input string) string {
	hash := sha1.Sum([]byte(input))

	return hex.EncodeToString(hash[:])
}
