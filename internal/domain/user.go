// Package domain defines core business logic and models for the application.
package domain

// User represents a user in the system.
//
// It encapsulates the core attributes of a user: username and password.
// This struct is used to represent user data across different layers of the application.
type User struct {
	username string
	password string
}
