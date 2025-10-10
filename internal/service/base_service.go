// Package service defines the base service interfaces for the OIDC application.
package service

// IBaseService defines the base service interface
type IBaseService interface {
	// GetDomainHelper returns the domain helper
	GetDomainHelper() interface{}

	// GetDomainMapper returns the domain mapper
	GetDomainMapper() interface{}
}

// ConvertInterface defines conversion methods between entities and DTOs
type ConvertInterface interface {
	// ConvertToEntity converts DTO to entity
	ConvertToEntity(dto interface{}) interface{}

	// ConvertToDTO converts entity to DTO
	ConvertToDTO(entity interface{}) interface{}
}