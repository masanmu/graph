package http

import (
	"net/http"

	"github.com/open-falcon/graph/index"
)

func configZtreeRoutes() {
	http.HandleFunc("/ztree/updateAll", func(w http.ResponseWriter, r *http.Request) {
		go index.UpdateZtreeAll()
		RenderDataJson(w, "ok")
	})
}
