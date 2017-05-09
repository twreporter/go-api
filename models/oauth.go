package models

// AuthenticatedResponse defines the user info fields
type AuthenticatedResponse struct {
	ID        uint   `json:"id"`
	Privilege int    `json:"privilege"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"eamil"`
	Jwt       string `json:"jwt"`
}
