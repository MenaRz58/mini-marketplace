package model

type Metadata struct {
	ID         string            `json:"id"`
	EntityID   string            `json:"entity_id"`
	EntityType string            `json:"entity_type"`
	Attributes map[string]string `json:"attributes"`
}
