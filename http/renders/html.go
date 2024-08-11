package renders

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/senither/dalamud-plugin-listing/metrics"
)

var (
	fileHash *string
)

func RenderHtml(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.HtmlRoute)

	fmt.Println(" - Rendering HTML")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	content, err := os.ReadFile("./views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	if fileHash == nil {
		hash := createSHA1Hash(string(content))
		fileHash = &hash
	}

	fmt.Fprint(w, strings.Replace(string(content), "@file-hash", *fileHash, 1))
}

func createSHA1Hash(input string) string {
	hash := sha1.Sum([]byte(input))

	return hex.EncodeToString(hash[:])
}
