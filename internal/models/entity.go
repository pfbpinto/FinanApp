package models

type Entity struct {
	EntityID       int    `json:"entity_id"`
	EntityName     string `json:"entity_name"`
	EntityType     string `json:"entity_type"`
	EntityCategory string `json:"entity_category"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
}
