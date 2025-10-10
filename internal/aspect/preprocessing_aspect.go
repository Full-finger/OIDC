// Package aspect provides cross-cutting concerns such as logging and preprocessing.
package aspect

import (
	"github.com/gin-gonic/gin"
)

// PreprocessingAspect implements data validation and preprocessing functionality
type PreprocessingAspect struct {
	enabled bool
}

// NewPreprocessingAspect creates a new PreprocessingAspect instance
func NewPreprocessingAspect() *PreprocessingAspect {
	return &PreprocessingAspect{
		enabled: true,
	}
}

// Before preprocesses data before controller execution
func (pa *PreprocessingAspect) Before(c *gin.Context) {
	if !pa.enabled {
		c.Next()
		return
	}

	// Here you can add data validation and preprocessing logic
	// For example:
	// 1. Validate request format
	// 2. Sanitize input data
	// 3. Transform data to match service requirements

	// For now, we just continue with the request
	c.Next()
}

// After handles post-processing after controller execution
func (pa *PreprocessingAspect) After(c *gin.Context) {
	if !pa.enabled {
		return
	}

	// Here you can add post-processing logic if needed
	// For example:
	// 1. Transform response data
	// 2. Add additional headers
}

// Handle preprocesses the complete request/response cycle
func (pa *PreprocessingAspect) Handle(c *gin.Context) {
	pa.Before(c)
	c.Next()
	pa.After(c)
}

// Enable enables preprocessing
func (pa *PreprocessingAspect) Enable() {
	pa.enabled = true
}

// Disable disables preprocessing
func (pa *PreprocessingAspect) Disable() {
	pa.enabled = false
}