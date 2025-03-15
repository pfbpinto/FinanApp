package handlers

import (
	"finanapp/internal/models"
	"html/template"
	"net/http"
)

// RenderTemplate loads and renders templates with the base layout
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	// Retrieve authentication and user data from the context
	authenticated := false
	user := ""

	// Access the context to get information about authentication and user
	if val := r.Context().Value("authenticated"); val != nil {
		authenticated = val.(bool)
	}

	if val := r.Context().Value("user"); val != nil {
		user = val.(models.User).EmailAddress
	}

	// Add authentication data to the data map
	finalData := map[string]interface{}{
		"Authenticated": authenticated,
		"User":          user,
	}

	// Add additional data passed to the template
	if data != nil {
		for key, value := range data.(map[string]interface{}) {
			finalData[key] = value
		}
	}

	// Load the template with the base layout
	templates, err := template.ParseFiles(
		"views/layout.html",   // Base layout
		"views/"+templateName, // Specific template
	)
	if err != nil {
		http.Error(w, "Error loading templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template with the final data
	err = templates.ExecuteTemplate(w, "layout.html", finalData)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}
