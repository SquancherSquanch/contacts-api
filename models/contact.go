package models

// Entries data model holding all contacts
type Entries struct {
	// Contacts ...
	Contacts []Contact `json:"contacts"`
}

// Contact a data model to store contact information
type Contact struct {
	// ID ...
	ID string `json:"id"`
	// FirstName ...
	FirstName string `json:"first_name"`
	// LastName ...
	LastName string `json:"last_name"`
	// Email ...
	Email string `json:"email"`
	// Phone ...
	Phone string `json:"phone"`
}

// HTTPErrorResponse response given when an error occurs from a http request
type HTTPErrorResponse struct {
	// Error ...
	Error string `json:"error,omitempty"`
}
