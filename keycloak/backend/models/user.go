package models

// User represents user information from Keycloak JWT token
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Sub   string `json:"sub"`
}

// UserResponse represents the response structure for user info API
type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
