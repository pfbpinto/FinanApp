package handlers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{}

	RenderTemplate(w, r, "home.html", data)
}
