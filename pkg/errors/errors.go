package errors

import "fmt"

// Common error types for the application

// ConfigError represents configuration-related errors
type ConfigError struct {
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("config error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("config error: %s", e.Message)
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, cause error) *ConfigError {
	return &ConfigError{Message: message, Cause: cause}
}

// TemplateError represents template-related errors
type TemplateError struct {
	Message string
	Cause   error
}

func (e *TemplateError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("template error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("template error: %s", e.Message)
}

// NewTemplateError creates a new template error
func NewTemplateError(message string, cause error) *TemplateError {
	return &TemplateError{Message: message, Cause: cause}
}

// ProcessorError represents processor-related errors
type ProcessorError struct {
	Message string
	Cause   error
}

func (e *ProcessorError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("processor error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("processor error: %s", e.Message)
}

// NewProcessorError creates a new processor error
func NewProcessorError(message string, cause error) *ProcessorError {
	return &ProcessorError{Message: message, Cause: cause}
}
