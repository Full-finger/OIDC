// Package handler defines the controller layer interfaces for the OIDC application.
package handler

import (
	"github.com/Full-finger/OIDC/internal/service"
)

// IBaseController defines the base controller interface
type IBaseController interface {
	// HealthCheck checks the health of the controller
	HealthCheck() error

	// GetVersion returns the version of the controller
	GetVersion() string
}

// BaseController provides common functionality for all controllers
type BaseController struct {
	version string
}

// NewBaseController creates a new BaseController instance
func NewBaseController() *BaseController {
	return &BaseController{
		version: "1.0.0",
	}
}

// HealthCheck checks the health of the controller
func (bc *BaseController) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the controller
func (bc *BaseController) GetVersion() string {
	return bc.version
}

// GetServiceHelper provides a helper function to get service
func GetServiceHelper(svc interface{}) service.IBaseService {
	if baseService, ok := svc.(service.IBaseService); ok {
		return baseService
	}
	return nil
}