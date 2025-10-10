// Package helper defines the base helper interfaces for the OIDC application.
package helper

// IBaseHelper defines the base helper interface
type IBaseHelper interface {
	// Validate validates the entity
	Validate(entity interface{}) error

	// Format formats the entity
	Format(entity interface{}) interface{}
}