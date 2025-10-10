// Package aspect provides cross-cutting concerns such as logging and preprocessing.
package aspect

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingAspect implements logging functionality for all controller methods
type LoggingAspect struct {
	enabled bool
}

// NewLoggingAspect creates a new LoggingAspect instance
func NewLoggingAspect() *LoggingAspect {
	return &LoggingAspect{
		enabled: true,
	}
}

// Before logs information before method execution
func (la *LoggingAspect) Before(c *gin.Context) {
	if !la.enabled {
		c.Next()
		return
	}

	startTime := time.Now()
	log.Printf("[REQUEST] %s %s - Start Time: %s", 
		c.Request.Method, 
		c.Request.URL.Path, 
		startTime.Format("2006-01-02 15:04:05"))

	// Store start time in context for After method
	c.Set("start_time", startTime)

	// Continue with the request
	c.Next()
}

// After logs information after method execution
func (la *LoggingAspect) After(c *gin.Context) {
	if !la.enabled {
		return
	}

	// Retrieve start time from context
	startTime, exists := c.Get("start_time")
	if !exists {
		return
	}

	duration := time.Since(startTime.(time.Time))
	log.Printf("[RESPONSE] %s %s - Status: %d - Duration: %v", 
		c.Request.Method, 
		c.Request.URL.Path, 
		c.Writer.Status(), 
		duration)
}

// Handle logs the complete request/response cycle
func (la *LoggingAspect) Handle(c *gin.Context) {
	la.Before(c)
	c.Next()
	la.After(c)
}

// Enable enables logging
func (la *LoggingAspect) Enable() {
	la.enabled = true
}

// Disable disables logging
func (la *LoggingAspect) Disable() {
	la.enabled = false
}

// LogError logs error messages
func (la *LoggingAspect) LogError(format string, args ...interface{}) {
	if la.enabled {
		log.Printf("[ERROR] "+format, args...)
	}
}

// LogInfo logs info messages
func (la *LoggingAspect) LogInfo(format string, args ...interface{}) {
	if la.enabled {
		log.Printf("[INFO] "+format, args...)
	}
}